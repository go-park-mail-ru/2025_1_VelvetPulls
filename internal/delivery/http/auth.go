package http

import (
	"net/http"

	apperrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/app_errors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	authpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/delivery/proto"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type authController struct {
	authClient    authpb.AuthServiceClient
	sessionClient authpb.SessionServiceClient
}

func NewAuthController(r *mux.Router, authClient authpb.AuthServiceClient, sessionClient authpb.SessionServiceClient) {
	controller := &authController{
		authClient:    authClient,
		sessionClient: sessionClient,
	}
	r.Handle("/register", http.HandlerFunc(controller.Register)).Methods(http.MethodPost)
	r.Handle("/login", http.HandlerFunc(controller.Login)).Methods(http.MethodPost)
	r.Handle("/logout", http.HandlerFunc(controller.Logout)).Methods(http.MethodDelete)
}

func (c *authController) Register(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())

	var creds model.RegisterCredentials
	if err := utils.ParseJSONRequest(r, &creds); err != nil {
		logger.Warn("Invalid request data", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid request data", false)
		return
	}

	resp, err := c.authClient.RegisterUser(r.Context(), &authpb.RegisterUserRequest{
		Username:        creds.Username,
		Password:        creds.Password,
		ConfirmPassword: creds.Password,
		Phone:           creds.Phone,
	})
	if err != nil {
		logger.Error("gRPC Register error", zap.Error(err))
		code, msg := apperrors.UnpackGrpcError(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SetSessionCookie(w, resp.GetSessionId())
	utils.SendJSONResponse(w, r, http.StatusCreated, "Registration successful", true)
}

func (c *authController) Login(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())

	var creds model.LoginCredentials
	if err := utils.ParseJSONRequest(r, &creds); err != nil {
		logger.Warn("Invalid request data", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid request data", false)
		return
	}

	resp, err := c.authClient.LoginUser(r.Context(), &authpb.LoginUserRequest{
		Username: creds.Username,
		Password: creds.Password,
	})
	if err != nil {
		logger.Error("gRPC Login error", zap.Error(err))
		code, msg := apperrors.UnpackGrpcError(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SetSessionCookie(w, resp.GetSessionId())
	utils.SendJSONResponse(w, r, http.StatusOK, "Login successful", true)
}

func (c *authController) Logout(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())

	sessionID, err := utils.GetSessionCookie(r)
	if err != nil {
		logger.Warn("Missing session cookie")
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Unauthorized", false)
		return
	}

	_, err = c.authClient.LogoutUser(r.Context(), &authpb.LogoutUserRequest{
		SessionId: sessionID,
	})
	if err != nil {
		logger.Error("gRPC Logout error", zap.Error(err))
		code, msg := apperrors.UnpackGrpcError(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.DeleteSessionCookie(w)
	utils.SendJSONResponse(w, r, http.StatusOK, "Logout successful", true)
}
