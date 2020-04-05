package storage

// CartItem ..
type CartItem struct {
	VariationID  int    `json:"variation_id"`
	СartSession  string `json:"cart_session"`
	SizeOptionID int    `json:"size_option_id"`
}

// CartItemResponse ..
type CartItemResponse struct {
	ID           int    `json:"variation_id"`
	Name         string `json:"name"`
	Brand        string `json:"brand"`
	Price        int    `json:"price"`
	PricePerItem int    `json:"price_per_item"`
	Size         string `json:"size"`
	Quantity     int    `json:"quantity"`
	SizeOptionID int    `json:"size_option_id"`
}
