package grpc

import (
	"context"
	"fmt"

	apperrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/app_errors"
	authpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/delivery/proto"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/usecase"
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
		Username:        req.GetUsername(),
		Password:        req.GetPassword(),
		ConfirmPassword: req.GetPassword(),
		Phone:           req.GetPhone(),
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
	fmt.Println(err)
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
	userID, err := c.sessionUsecase.CheckLogin(ctx, req.GetSessionId())
	if err != nil {
		return nil, apperrors.ConvertError(err)
	}
	return &authpb.CheckLoginResponse{UserId: userID}, nil
}
