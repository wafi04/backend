package request

type CartRequest struct {
	VariantID string  `json:"variant_id"`
	Size      string  `json:"size"`
	Quantity  int64   `json:"quantity"`
	UserID    string  `json:"user_id"`
	Total     float64 `json:"total"`
}

type ReqRemoveCartByID struct {
	CartItemID string `json:"cart_item_id"`
	UserID     string `json:"user_id"`
}

type ClearCart struct {
	CartID string `json:"cart_id"`
}

type UpdateQuantity struct {
	CartItemID string `json:"cart_item_id"`
	Size       string `json:"size"`
	Quantity   int64  `json:"quantity"`
}
