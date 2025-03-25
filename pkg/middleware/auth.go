package middleware

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
)

func AuthMiddleware(sessionUC usecase.ISessionUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := utils.GetSessionCookie(r)
			if err != nil {
				utils.SendJSONResponse(w, http.StatusBadRequest, "Unauthorized", false)
				return
			}

			userID, err := sessionUC.CheckLogin(r.Context(), token)
			if err != nil {
				utils.SendJSONResponse(w, http.StatusUnauthorized, "Invalid session", false)
				return
			}

			ctxWithUID := context.WithValue(r.Context(), utils.USER_ID_KEY, userID)

			next.ServeHTTP(w, r.WithContext(ctxWithUID))
		})
	}
}
