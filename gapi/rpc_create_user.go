package gapi

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	db "github.com/zjr71163356/simplebank/db/sqlc"
	"github.com/zjr71163356/simplebank/pb"
	"github.com/zjr71163356/simplebank/utils"
	"github.com/zjr71163356/simplebank/val"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if err := ValidateCreateUserRequest(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	HashedPassword, err := utils.HashPassWord(req.Password)
	if err != nil {

		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: HashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	})

	if err != nil {
		//foreign_key_violation
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user already exists: %v", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "fail to create user: %v ", err)
	}
	rsp := &pb.CreateUserResponse{
		User: &pb.User{
			Username:          user.Username,
			FullName:          user.FullName,
			Email:             user.Email,
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
			CreatedAt:         timestamppb.New(user.CreatedAt),
		},
	}

	return rsp, nil

}

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	if err := ValidateLoginUserRequest(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %s", err)
	}

	err = utils.MatchPassWord(user.HashedPassword, req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect password")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %s", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %s", err)
	}
	metadata := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.Id,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err)
	}

	rsp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return rsp, nil
}

func ValidateCreateUserRequest(req *pb.CreateUserRequest) error {
	if err := val.ValidateUsername(req.Username); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid username: %s", err)
	}
	if err := val.ValidatePassword(req.Password); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid password: %s", err)
	}
	if err := val.ValidateFullName(req.FullName); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid full name: %s", err)
	}
	if err := val.ValidateEmail(req.Email); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid email: %s", err)
	}
	return nil
}

func ValidateLoginUserRequest(req *pb.LoginUserRequest) error {

	if err := val.ValidateUsername(req.Username); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid username: %s", err)
	}
	if err := val.ValidatePassword(req.Password); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid password: %s", err)
	}
	return nil
}
