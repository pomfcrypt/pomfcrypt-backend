package model

type FileResponse struct {
	Hash       [32]byte `json:"hash"`
	Filename   string   `json:"filename"`
	UploadedAt int64    `json:"uploaded_at"`
}
