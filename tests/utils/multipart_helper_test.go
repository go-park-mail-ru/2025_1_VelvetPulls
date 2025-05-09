package utils_test

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type multipartFileMock struct {
	*bytes.Reader
}

func newMultipartFileMock(data []byte) multipart.File {
	return &multipartFileMock{Reader: bytes.NewReader(data)}
}

func (m *multipartFileMock) Close() error {
	return nil
}

func (m *multipartFileMock) Seek(offset int64, whence int) (int64, error) {
	return m.Reader.Seek(offset, whence)
}

func TestSavePhotoInvalid(t *testing.T) {
	file := newMultipartFileMock([]byte("not an image"))
	_, err := utils.SavePhoto(file, "avatars")
	assert.ErrorIs(t, err, utils.ErrNotImage)
}

func TestIsImageFile(t *testing.T) {
	img := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG header
	text := []byte("this is not an image")

	imgFile := newMultipartFileMock(img)
	textFile := newMultipartFileMock(text)

	assert.True(t, utils.IsImageFile(imgFile))
	assert.False(t, utils.IsImageFile(textFile))
}

func TestSavePhoto(t *testing.T) {
	testDir := filepath.Join(config.UPLOAD_DIR, "test")
	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	tests := []struct {
		name        string
		fileContent []byte
		contentType string
		folderName  string
		wantErr     error
	}{
		{
			name:        "Valid image",
			fileContent: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, // PNG header
			contentType: "image/png",
			folderName:  "test",
			wantErr:     nil,
		},
		{
			name:        "Not an image",
			fileContent: []byte("not an image"),
			contentType: "text/plain",
			folderName:  "test",
			wantErr:     utils.ErrNotImage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("file", "test.png")
			require.NoError(t, err)

			_, err = part.Write(tt.fileContent)
			require.NoError(t, err)
			writer.Close()

			req := httptest.NewRequest("POST", "/", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			file, _, err := req.FormFile("file")
			require.NoError(t, err)
			defer file.Close()

			path, err := utils.SavePhoto(file, tt.folderName)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, path)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, path)
				assert.FileExists(t, path)
			}
		})
	}
}

func TestRewritePhoto(t *testing.T) {
	testDir := filepath.Join(config.UPLOAD_DIR, "test")
	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFile := filepath.Join(testDir, "test.png")
	err = os.WriteFile(testFile, []byte("old content"), 0644)
	require.NoError(t, err)

	t.Run("Success rewrite", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.png")
		require.NoError(t, err)

		_, err = part.Write([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) // valid PNG
		require.NoError(t, err)
		writer.Close()

		req := httptest.NewRequest("POST", "/", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		file, _, err := req.FormFile("file")
		require.NoError(t, err)
		defer file.Close()

		err = utils.RewritePhoto(file, testFile)
		assert.NoError(t, err)

		content, err := os.ReadFile(testFile)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, content)
	})

}

func TestRemovePhoto(t *testing.T) {
	testDir := filepath.Join(config.UPLOAD_DIR, "test")
	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFile := filepath.Join(testDir, "test.png")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	t.Run("Success remove", func(t *testing.T) {
		err := utils.RemovePhoto(testFile)
		assert.NoError(t, err)
		assert.NoFileExists(t, testFile)
	})

	t.Run("File not exists", func(t *testing.T) {
		err := utils.RemovePhoto("nonexistent/path/test.png")
		assert.ErrorIs(t, err, utils.ErrDeletingImage)
	})
}
