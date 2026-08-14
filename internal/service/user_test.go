package service_test

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	userv1 "github.com/example/grpc-user-service/gen/user/v1"
	"github.com/example/grpc-user-service/internal/service"
	"github.com/example/grpc-user-service/internal/store"
)

const bufSize = 1024 * 1024

func startTestServer(t *testing.T) (userv1.UserServiceClient, func()) {
	t.Helper()

	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	userv1.RegisterUserServiceServer(srv, service.NewUserServer(store.NewMemoryStore()))

	go func() {
		_ = srv.Serve(lis)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	cancel()
	if err != nil {
		srv.Stop()
		t.Fatalf("dial bufnet: %v", err)
	}

	cleanup := func() {
		_ = conn.Close()
		srv.Stop()
		_ = lis.Close()
	}
	return userv1.NewUserServiceClient(conn), cleanup
}

func TestUserServiceCRUD(t *testing.T) {
	client, cleanup := startTestServer(t)
	defer cleanup()

	ctx := context.Background()

	created, err := client.CreateUser(ctx, &userv1.CreateUserRequest{
		Name:  "Bob",
		Email: "bob@example.com",
	})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if created.User.Id == "" {
		t.Fatal("expected non-empty id")
	}

	got, err := client.GetUser(ctx, &userv1.GetUserRequest{Id: created.User.Id})
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got.User.Email != "bob@example.com" {
		t.Fatalf("unexpected email: %s", got.User.Email)
	}

	updated, err := client.UpdateUser(ctx, &userv1.UpdateUserRequest{
		Id:    created.User.Id,
		Name:  "Bobby",
		Email: "bobby@example.com",
	})
	if err != nil {
		t.Fatalf("UpdateUser: %v", err)
	}
	if updated.User.Name != "Bobby" {
		t.Fatalf("unexpected name: %s", updated.User.Name)
	}

	listed, err := client.ListUsers(ctx, &userv1.ListUsersRequest{PageSize: 10})
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if len(listed.Users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(listed.Users))
	}

	if _, err := client.DeleteUser(ctx, &userv1.DeleteUserRequest{Id: created.User.Id}); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}

	_, err = client.GetUser(ctx, &userv1.GetUserRequest{Id: created.User.Id})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("expected NotFound, got %v", err)
	}
}

func TestUserServiceErrors(t *testing.T) {
	client, cleanup := startTestServer(t)
	defer cleanup()
	ctx := context.Background()

	_, err := client.CreateUser(ctx, &userv1.CreateUserRequest{Name: "", Email: "x@example.com"})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", err)
	}

	_, err = client.CreateUser(ctx, &userv1.CreateUserRequest{Name: "A", Email: "a@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.CreateUser(ctx, &userv1.CreateUserRequest{Name: "B", Email: "a@example.com"})
	if status.Code(err) != codes.AlreadyExists {
		t.Fatalf("expected AlreadyExists, got %v", err)
	}

	_, err = client.GetUser(ctx, &userv1.GetUserRequest{Id: "missing"})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("expected NotFound, got %v", err)
	}
}

func TestUserServiceListPagination(t *testing.T) {
	client, cleanup := startTestServer(t)
	defer cleanup()
	ctx := context.Background()

	emails := []string{"u1@example.com", "u2@example.com", "u3@example.com"}
	for _, email := range emails {
		_, err := client.CreateUser(ctx, &userv1.CreateUserRequest{
			Name:  "U",
			Email: email,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	page1, err := client.ListUsers(ctx, &userv1.ListUsersRequest{PageSize: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(page1.Users) != 2 || page1.NextPageToken == "" {
		t.Fatalf("unexpected page1: users=%d token=%q", len(page1.Users), page1.NextPageToken)
	}

	page2, err := client.ListUsers(ctx, &userv1.ListUsersRequest{
		PageSize:  2,
		PageToken: page1.NextPageToken,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(page2.Users) != 1 {
		t.Fatalf("expected 1 user on page2, got %d", len(page2.Users))
	}
}
