package middleware

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	generatedAuth "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/proto"
	"github.com/google/uuid"
)

func AuthMiddleware(sessionClient generatedAuth.SessionServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := utils.GetSessionCookie(r)
			if err != nil {
				utils.SendJSONResponse(w, r, http.StatusUnauthorized, "Unauthorized", false)
				return
			}

			resp, err := sessionClient.CheckLogin(r.Context(), &generatedAuth.CheckLoginRequest{
				SessionId: token,
			})
			if err != nil {
				utils.DeleteSessionCookie(w)
				utils.SendJSONResponse(w, r, http.StatusUnauthorized, "Invalid session", false)
				return
			}
			userID, err := uuid.Parse(resp.UserId)
			if err != nil {
				utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid user ID", false)
				return
			}

			ctxWithUID := context.WithValue(r.Context(), utils.USER_ID_KEY, userID)
			next.ServeHTTP(w, r.WithContext(ctxWithUID))
		})
	}
}

func AuthMiddlewareWS(sessionClient generatedAuth.SessionServiceClient) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token, err := utils.GetSessionCookie(r)
			if err != nil {
				utils.SendJSONResponse(w, r, http.StatusUnauthorized, "Unauthorized", false)
				return
			}

			resp, err := sessionClient.CheckLogin(r.Context(), &generatedAuth.CheckLoginRequest{
				SessionId: token,
			})
			if err != nil {
				utils.DeleteSessionCookie(w)
				utils.SendJSONResponse(w, r, http.StatusUnauthorized, "Invalid session", false)
				return
			}

			userID, err := uuid.Parse(resp.UserId)
			if err != nil {
				utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid user ID", false)
				return
			}

			ctx := context.WithValue(r.Context(), utils.USER_ID_KEY, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
