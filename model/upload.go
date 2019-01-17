package model

type FileResponse struct {
	Hash       string `json:"hash"`
	Filename   string `json:"filename"`
	UploadedAt int64  `json:"uploaded_at"`
}
