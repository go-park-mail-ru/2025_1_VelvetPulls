package http

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	authpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/proto"
	"github.com/gorilla/mux"
)

type filesController struct {
	sessionClient authpb.SessionServiceClient
	filesUsecase  usecase.IFilesUsecase
}

func NewFilesController(r *mux.Router, sessionClient authpb.SessionServiceClient, filesUsecase usecase.IFilesUsecase) {
	controller := &filesController{
		sessionClient: sessionClient,
		filesUsecase:  filesUsecase,
	}

	r.Handle("/files/{file_id}", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.GetFile))).Methods(http.MethodGet)
}

func (c *filesController) GetFile(w http.ResponseWriter, r *http.Request) {

}

type File struct {
	*os.File
}

func (f *File) Close() error {
	return f.File.Close()
}

func newFileHeader(filePath string) (*multipart.FileHeader, error) {
	// Получаем информацию о файле
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	// Заполняем заголовок
	header := &multipart.FileHeader{
		Filename: fileInfo.Name(),
		Size:     fileInfo.Size(),
		Header:   make(textproto.MIMEHeader),
	}
	header.Header.Set("Content-Type", "image/webp")

	fmt.Print("sticker header: ", header)

	return header, nil
}

func getMultipartFile(filePath string) (multipart.File, *multipart.FileHeader, error) {
	// Открываем файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Errorf("файл не существует: %s", filePath)
		return nil, nil, fmt.Errorf("файл не существует: %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Errorf("sticker error open: %v", err)
		return nil, nil, err
	}

	fmt.Println("handler sticker: ", file)

	fileInfo, err := file.Stat()
	if err != nil {
		file.Close() // Закрываем файл при ошибке
		fmt.Errorf("нет инфы файла: %s", filePath)
		return nil, nil, err
	}
	if fileInfo.Size() == 0 {
		file.Close() // Закрываем файл при отсутствии данных
		fmt.Errorf("файл пустой: %s", filePath)
		return nil, nil, fmt.Errorf("файл пустой: %s", filePath)
	}

	// Получаем заголовок
	header, err := newFileHeader(filePath)
	if err != nil {
		file.Close() // Закрываем файл при ошибке
		fmt.Errorf("sticker header error open: %v", err)
		return nil, nil, err
	}

	return &File{File: file}, header, nil
}
