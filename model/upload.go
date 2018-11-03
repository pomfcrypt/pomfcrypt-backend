package model

type FileResponse struct {
	Message    string `json:"message"`
	Hash       string `json:"hash"`
	Name       string `json:"name"`
	UploadedAt int64  `json:"uploaded_at"`
}
