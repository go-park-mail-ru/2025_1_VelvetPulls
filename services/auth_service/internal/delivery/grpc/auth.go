package grpc

import (
	"context"

	apperrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/internal/app_errors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/internal/usecase"
	authpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type authController struct {
	authpb.UnimplementedAuthServiceServer
	authpb.UnimplementedSessionServiceServer

	authUsecase    usecase.IAuthUsecase
	sessionUsecase usecase.ISessionUsecase
}

func NewAuthController(grpcServer *grpc.Server, authUsecase usecase.IAuthUsecase, sessionUsecase usecase.ISessionUsecase) {
	controller := &authController{
		authUsecase:    authUsecase,
		sessionUsecase: sessionUsecase,
	}
	authpb.RegisterAuthServiceServer(grpcServer, controller)
	authpb.RegisterSessionServiceServer(grpcServer, controller)
}

func (c *authController) RegisterUser(ctx context.Context, req *authpb.RegisterUserRequest) (*authpb.RegisterUserResponse, error) {
	sessionID, err := c.authUsecase.RegisterUser(ctx, model.RegisterCredentials{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Name:     req.GetName(),
	})
	if err != nil {
		return nil, apperrors.ConvertError(err)
	}

	return &authpb.RegisterUserResponse{SessionId: sessionID}, nil
}

func (c *authController) LoginUser(ctx context.Context, req *authpb.LoginUserRequest) (*authpb.LoginUserResponse, error) {
	sessionID, err := c.authUsecase.LoginUser(ctx, model.LoginCredentials{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, apperrors.ConvertError(err)
	}

	return &authpb.LoginUserResponse{SessionId: sessionID}, nil
}

func (c *authController) LogoutUser(ctx context.Context, req *authpb.LogoutUserRequest) (*emptypb.Empty, error) {
	err := c.authUsecase.LogoutUser(ctx, req.GetSessionId())
	if err != nil {
		return nil, apperrors.ConvertError(err)
	}
	return &emptypb.Empty{}, nil
}

func (c *authController) CheckLogin(ctx context.Context, req *authpb.CheckLoginRequest) (*authpb.CheckLoginResponse, error) {
	user, err := c.sessionUsecase.CheckLogin(ctx, req.GetSessionId())
	if err != nil {
		return nil, apperrors.ConvertError(err)
	}

	var avatar string
	if user.AvatarPath != nil {
		avatar = *user.AvatarPath
	}

	return &authpb.CheckLoginResponse{
		UserId:   user.ID.String(),
		Username: user.Username,
		Name:     user.Name,
		Avatar:   avatar,
	}, nil
}
