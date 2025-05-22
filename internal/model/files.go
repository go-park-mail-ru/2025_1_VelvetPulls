package model

import "github.com/google/uuid"

type UploadFileResponse struct {
	FileID string `json:"file_id"`
}

type FileMetaData struct {
	Filename    string
	ContentType string
	FileSize    int64
}

type Payload struct {
	URL      string
	Filename string
	Size     int64
}

type GetStickerPackResponse struct {
	Photo string   `json:"photo" valid:"-"`
	URLs  []string `json:"stickers" valid:"-"`
}

type StickerPack struct {
	Photo  string    `json:"photo" valid:"-"`
	Name   string    `json:"name"`
	PackID uuid.UUID `json:"id" valid:"-"`
}

type StickerPacks struct {
	Packs []StickerPack `json:"packs" valid:"-"`
}
