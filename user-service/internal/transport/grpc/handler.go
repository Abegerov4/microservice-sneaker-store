package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"sneaker-store/user-service/internal/model"
	"sneaker-store/user-service/internal/usecase"
	pb "sneaker-store/user-service/proto"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	uc *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

func toProto(u *model.User) *pb.UserResponse {
	return &pb.UserResponse{
		Id:        u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		Phone:     u.Phone,
		Active:    u.Active,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.String(),
		UpdatedAt: u.UpdatedAt.String(),
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	role := req.Role
	if role == "" {
		role = model.RoleUser
	}

	u, err := h.uc.CreateWithRole(ctx, req.Email, req.Password, req.FullName, req.Phone, role)
	if err != nil {
		if err.Error() == "email already in use" {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "create user: %v", err)
	}
	return toProto(u), nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	u, err := h.uc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}
	return toProto(u), nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	u, err := h.uc.Update(ctx, req.Id, req.FullName, req.Phone)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update user: %v", err)
	}
	return toProto(u), nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if err := h.uc.Delete(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.NotFound, "delete user: %v", err)
	}
	return &pb.DeleteUserResponse{Success: true}, nil
}

func (h *UserHandler) AuthenticateUser(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}
	u, err := h.uc.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		return &pb.AuthenticateResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.AuthenticateResponse{Success: true, UserId: u.ID, Role: u.Role, Message: "authenticated"}, nil
}

func (h *UserHandler) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	u, err := h.uc.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}
	return toProto(u), nil
}

func (h *UserHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	if req.Id == "" || req.OldPassword == "" || req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "id, old_password and new_password are required")
	}
	if err := h.uc.ChangePassword(ctx, req.Id, req.OldPassword, req.NewPassword); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "change password: %v", err)
	}
	return &pb.ChangePasswordResponse{Success: true}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, total, err := h.uc.List(ctx, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list users: %v", err)
	}
	resp := &pb.ListUsersResponse{Total: int32(total)}
	for _, u := range users {
		resp.Users = append(resp.Users, toProto(u))
	}
	return resp, nil
}

func (h *UserHandler) SearchUsers(ctx context.Context, req *pb.SearchUsersRequest) (*pb.ListUsersResponse, error) {
	if req.Query == "" {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}
	users, err := h.uc.Search(ctx, req.Query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "search users: %v", err)
	}
	resp := &pb.ListUsersResponse{Total: int32(len(users))}
	for _, u := range users {
		resp.Users = append(resp.Users, toProto(u))
	}
	return resp, nil
}

func (h *UserHandler) GetUserStats(ctx context.Context, _ *pb.GetUserStatsRequest) (*pb.GetUserStatsResponse, error) {
	stats, err := h.uc.GetStats(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get user stats: %v", err)
	}
	return &pb.GetUserStatsResponse{
		TotalUsers:  int32(stats.TotalUsers),
		ActiveUsers: int32(stats.ActiveUsers),
	}, nil
}

func (h *UserHandler) UpdateUserStatus(ctx context.Context, req *pb.UpdateUserStatusRequest) (*pb.UserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	u, err := h.uc.UpdateStatus(ctx, req.Id, req.Active)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update user status: %v", err)
	}
	return toProto(u), nil
}

func (h *UserHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	if req.Id == "" || req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "id and new_password are required")
	}
	if err := h.uc.ResetPassword(ctx, req.Id, req.NewPassword); err != nil {
		return nil, status.Errorf(codes.Internal, "reset password: %v", err)
	}
	return &pb.ResetPasswordResponse{Success: true}, nil
}
