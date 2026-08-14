package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userv1 "github.com/example/grpc-user-service/gen/user/v1"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:50051", "server address")
	action := flag.String("action", "demo", "action: demo|create|get|list|update|delete")
	id := flag.String("id", "", "user id")
	name := flag.String("name", "Alice", "user name")
	email := flag.String("email", "alice@example.com", "user email")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, *addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	client := userv1.NewUserServiceClient(conn)

	switch *action {
	case "demo":
		runDemo(ctx, client)
	case "create":
		resp, err := client.CreateUser(ctx, &userv1.CreateUserRequest{Name: *name, Email: *email})
		must(err)
		printUser(resp.User)
	case "get":
		resp, err := client.GetUser(ctx, &userv1.GetUserRequest{Id: *id})
		must(err)
		printUser(resp.User)
	case "list":
		resp, err := client.ListUsers(ctx, &userv1.ListUsersRequest{PageSize: 20})
		must(err)
		for _, u := range resp.Users {
			printUser(u)
		}
	case "update":
		resp, err := client.UpdateUser(ctx, &userv1.UpdateUserRequest{Id: *id, Name: *name, Email: *email})
		must(err)
		printUser(resp.User)
	case "delete":
		_, err := client.DeleteUser(ctx, &userv1.DeleteUserRequest{Id: *id})
		must(err)
		fmt.Println("deleted", *id)
	default:
		fmt.Fprintf(os.Stderr, "unknown action %q\n", *action)
		os.Exit(2)
	}
}

func runDemo(ctx context.Context, client userv1.UserServiceClient) {
	created, err := client.CreateUser(ctx, &userv1.CreateUserRequest{
		Name:  "Alice",
		Email: "alice@example.com",
	})
	must(err)
	fmt.Println("created:")
	printUser(created.User)

	got, err := client.GetUser(ctx, &userv1.GetUserRequest{Id: created.User.Id})
	must(err)
	fmt.Println("got:")
	printUser(got.User)

	listed, err := client.ListUsers(ctx, &userv1.ListUsersRequest{PageSize: 10})
	must(err)
	fmt.Printf("listed %d user(s)\n", len(listed.Users))
}

func printUser(u *userv1.User) {
	fmt.Printf("  id=%s name=%s email=%s\n", u.Id, u.Name, u.Email)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
