package service

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userv1 "github.com/example/grpc-user-service/gen/user/v1"
	"github.com/example/grpc-user-service/internal/store"
)

// UserServer implements userv1.UserServiceServer.
type UserServer struct {
	userv1.UnimplementedUserServiceServer
	store *store.MemoryStore
}

func NewUserServer(s *store.MemoryStore) *UserServer {
	return &UserServer{store: s}
}

func (s *UserServer) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	u, err := s.store.Create(req.GetName(), req.GetEmail())
	if err != nil {
		return nil, mapError(err)
	}
	return &userv1.CreateUserResponse{User: toProto(u)}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	u, err := s.store.Get(req.GetId())
	if err != nil {
		return nil, mapError(err)
	}
	return &userv1.GetUserResponse{User: toProto(u)}, nil
}

func (s *UserServer) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	users, next, err := s.store.List(int(req.GetPageSize()), req.GetPageToken())
	if err != nil {
		return nil, mapError(err)
	}
	out := make([]*userv1.User, 0, len(users))
	for _, u := range users {
		out = append(out, toProto(u))
	}
	return &userv1.ListUsersResponse{Users: out, NextPageToken: next}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	u, err := s.store.Update(req.GetId(), req.GetName(), req.GetEmail())
	if err != nil {
		return nil, mapError(err)
	}
	return &userv1.UpdateUserResponse{User: toProto(u)}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := s.store.Delete(req.GetId()); err != nil {
		return nil, mapError(err)
	}
	return &userv1.DeleteUserResponse{}, nil
}

func toProto(u *store.User) *userv1.User {
	return &userv1.User{
		Id:            u.ID,
		Name:          u.Name,
		Email:         u.Email,
		CreatedAtUnix: u.CreatedAt.Unix(),
		UpdatedAtUnix: u.UpdatedAt.Unix(),
	}
}

func mapError(err error) error {
	switch {
	case errors.Is(err, store.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, store.ErrEmailExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, store.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Errorf(codes.Internal, "internal error: %v", err)
	}
}
