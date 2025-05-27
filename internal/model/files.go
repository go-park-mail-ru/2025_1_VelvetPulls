//go:generate easyjson -all files.go
package model

import "github.com/google/uuid"

//easyjson:json
type UploadFileResponse struct {
	FileID string `json:"file_id"`
}

type FileMetaData struct {
	Filename    string
	ContentType string
	FileSize    int64
}

//easyjson:json
type Payload struct {
	URL         string
	Filename    string
	ContentType string
	Size        int64
}

//easyjson:json
type GetStickerPackResponse struct {
	Photo string   `json:"photo" valid:"-"`
	URLs  []string `json:"stickers" valid:"-"`
}

//easyjson:json
type StickerPack struct {
	Photo  string    `json:"photo" valid:"-"`
	Name   string    `json:"name"`
	PackID uuid.UUID `json:"id" valid:"-"`
}

type StickerPacks struct {
	Packs []StickerPack `json:"packs" valid:"-"`
}
