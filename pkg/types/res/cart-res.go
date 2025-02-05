package response

type CartResponse struct {
	VariantID string `json:"variant_id"`
	Size      string `json:"size"`
	Quantity  int64  `json:"quantity"`
}
type ResRemoveCartItem struct {
	Success bool `json:"success"`
}
