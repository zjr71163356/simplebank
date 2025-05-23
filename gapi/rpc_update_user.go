package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/zjr71163356/simplebank/db/sqlc"
	"github.com/zjr71163356/simplebank/pb"
	"github.com/zjr71163356/simplebank/utils"
	"github.com/zjr71163356/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if violations := ValidateUpdateUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	payload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to authorize user: %v", err)
	}

	if payload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "user does not have permission to update user ")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.GetFullName() != "",
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.GetEmail() != "",
		},
	}

	if req.GetPassword() != "" {
		HashedPassword, err := utils.HashPassWord(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
		}
		arg.HashedPassword = sql.NullString{
			String: HashedPassword,
			Valid:  true,
		}
		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	updatedUser, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: &pb.User{
			Username:          updatedUser.Username,
			FullName:          updatedUser.FullName,
			Email:             updatedUser.Email,
			PasswordChangedAt: timestamppb.New(updatedUser.PasswordChangedAt),
			CreatedAt:         timestamppb.New(updatedUser.CreatedAt),
		},
	}

	return rsp, nil
}

func ValidateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if req.GetPassword() != "" {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	if req.GetFullName() != "" {
		if err := val.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}

	if req.GetEmail() != "" {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	return violations
}
