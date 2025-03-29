package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

type uploadsController struct {
}

func NewUploadsController(r *mux.Router) {
	controller := &uploadsController{}

	r.HandleFunc("/uploads/{folder}/{name}", controller.GetImage).Methods(http.MethodGet)
}

// GetImage отправляет клиенту файл из папки загрузок.
//
// @Summary Получение загруженного файла
// @Description Возвращает файл из указанной папки на сервере
// @Tags Uploads
// @Produce octet-stream
// @Param folder path string true "Название папки"
// @Param name path string true "Имя файла"
// @Success 200 {file} binary
// @Failure 404 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /uploads/{folder}/{name} [get]
func (d *uploadsController) GetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folder := vars["folder"]
	name := vars["name"]

	imagePath := "./uploads/" + folder + "/" + name
	http.ServeFile(w, r, imagePath)
}
