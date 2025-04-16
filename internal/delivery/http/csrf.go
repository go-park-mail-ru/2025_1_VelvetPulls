package http

import (
	"net/http"

	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/gorilla/csrf"
)

func GetCSRFTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := csrf.Token(r)
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{
		"csrf_token": token,
	}, true)
}
