package auth

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase/auth"
	"github.com/gorilla/mux"
)

type authController struct {
	authUsecase auth.AuthUsecaseInterface
}

func NewAuthController(r *mux.Router, authUsecase auth.AuthUsecaseInterface) {
	controller := &authController{
		authUsecase: authUsecase,
	}

	r.HandleFunc("/register/", controller.Register).Methods(http.MethodPost)
	r.HandleFunc("/login/", controller.Login).Methods(http.MethodPost)
	// r.HandleFunc("/logout/", controller.Logout).Methods(http.MethodPost) TODO: сделать logout
}
