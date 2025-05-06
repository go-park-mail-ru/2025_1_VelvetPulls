package grpc_test

// func TestRegisterUser_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	authMock := mocks.NewMockIAuthUsecase(ctrl)
// 	sessionMock := mocks.NewMockISessionUsecase(ctrl)

// 	// Create controller directly without server
// 	c := &grpc.AuthController{
// 		authUsecase:    authMock,
// 		sessionUsecase: sessionMock,
// 	}

// 	req := &authpb.RegisterUserRequest{
// 		Username: "testuser",
// 		Password: "password",
// 		Phone:    "1234567890",
// 	}

// 	authMock.EXPECT().RegisterUser(gomock.Any(), model.RegisterCredentials{
// 		Username:        "testuser",
// 		Password:        "password",
// 		ConfirmPassword: "password",
// 		Phone:           "1234567890",
// 	}).Return("session123", nil)

// 	resp, err := c.RegisterUser(context.Background(), req)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "session123", resp.GetSessionId())
// }

// func TestRegisterUser_Error(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	authMock := mocks.NewMockIAuthUsecase(ctrl)
// 	sessionMock := mocks.NewMockISessionUsecase(ctrl)

// 	c := &grpc.AuthController{
// 		authUsecase:    authMock,
// 		sessionUsecase: sessionMock,
// 	}

// 	req := &authpb.RegisterUserRequest{
// 		Username: "testuser",
// 		Password: "password",
// 		Phone:    "1234567890",
// 	}

// 	expectedErr := errors.New("some error")
// 	authMock.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return("", expectedErr)

// 	resp, err := c.RegisterUser(context.Background(), req)
// 	assert.Nil(t, resp)
// 	assert.Error(t, err)

// 	grpcErr, ok := status.FromError(err)
// 	assert.True(t, ok)
// 	assert.Equal(t, codes.Internal, grpcErr.Code())
// }

// func TestLoginUser_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	authMock := mocks.NewMockIAuthUsecase(ctrl)
// 	sessionMock := mocks.NewMockISessionUsecase(ctrl)

// 	c := &grpc.AuthController{
// 		authUsecase:    authMock,
// 		sessionUsecase: sessionMock,
// 	}

// 	req := &authpb.LoginUserRequest{
// 		Username: "testuser",
// 		Password: "password",
// 	}

// 	authMock.EXPECT().LoginUser(gomock.Any(), model.LoginCredentials{
// 		Username: "testuser",
// 		Password: "password",
// 	}).Return("session123", nil)

// 	resp, err := c.LoginUser(context.Background(), req)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "session123", resp.GetSessionId())
// }

// func TestLogoutUser_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	authMock := mocks.NewMockIAuthUsecase(ctrl)
// 	sessionMock := mocks.NewMockISessionUsecase(ctrl)

// 	c := &grpc.AuthController{
// 		authUsecase:    authMock,
// 		sessionUsecase: sessionMock,
// 	}

// 	req := &authpb.LogoutUserRequest{SessionId: "session123"}
// 	authMock.EXPECT().LogoutUser(gomock.Any(), "session123").Return(nil)

// 	resp, err := c.LogoutUser(context.Background(), req)
// 	assert.NoError(t, err)
// 	assert.Equal(t, &emptypb.Empty{}, resp)
// }

// func TestLogoutUser_Error(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	authMock := mocks.NewMockIAuthUsecase(ctrl)
// 	sessionMock := mocks.NewMockISessionUsecase(ctrl)

// 	c := &grpc.AuthController{
// 		authUsecase:    authMock,
// 		sessionUsecase: sessionMock,
// 	}

// 	req := &authpb.LogoutUserRequest{SessionId: "session123"}
// 	expectedErr := errors.New("logout error")
// 	authMock.EXPECT().LogoutUser(gomock.Any(), "session123").Return(expectedErr)

// 	resp, err := c.LogoutUser(context.Background(), req)
// 	assert.Nil(t, resp)
// 	assert.Error(t, err)

// 	grpcErr, ok := status.FromError(err)
// 	assert.True(t, ok)
// 	assert.Equal(t, codes.Internal, grpcErr.Code())
// }

// func TestCheckLogin_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	authMock := mocks.NewMockIAuthUsecase(ctrl)
// 	sessionMock := mocks.NewMockISessionUsecase(ctrl)

// 	c := &grpc.AuthController{
// 		authUsecase:    authMock,
// 		sessionUsecase: sessionMock,
// 	}

// 	req := &authpb.CheckLoginRequest{SessionId: "session123"}
// 	sessionMock.EXPECT().CheckLogin(gomock.Any(), "session123").Return("user123", nil)

// 	resp, err := c.CheckLogin(context.Background(), req)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "user123", resp.GetUserId())
// }

// func TestCheckLogin_Error(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	authMock := mocks.NewMockIAuthUsecase(ctrl)
// 	sessionMock := mocks.NewMockISessionUsecase(ctrl)

// 	c := &grpc.AuthController{
// 		authUsecase:    authMock,
// 		sessionUsecase: sessionMock,
// 	}

// 	req := &authpb.CheckLoginRequest{SessionId: "invalid"}
// 	expectedErr := errors.New("session not found")
// 	sessionMock.EXPECT().CheckLogin(gomock.Any(), "invalid").Return("", expectedErr)

// 	resp, err := c.CheckLogin(context.Background(), req)
// 	assert.Nil(t, resp)
// 	assert.Error(t, err)

// 	grpcErr, ok := status.FromError(err)
// 	assert.True(t, ok)
// 	assert.Equal(t, codes.Internal, grpcErr.Code())
// }
