package response

type FileUploadResponse struct {
	URL      string `json:"url"`
	PublicID string `json:"public_id"`
	Error    string `json:"error,omitempty"`
}
