package request

import "mime/multipart"

type FileUploadRequest struct {
	FileData multipart.File `json:"file_data"`
	FileName string         `json:"file_name"`
	Folder   string         `json:"folder"`
	PublicID string         `json:"public_id"`
	FileType FileType       `json:"file_type"`
}

type FileType int32

const (
	IMAGE    FileType = 0
	VIDEO    FileType = 1
	DOCUMENT FileType = 2
)
