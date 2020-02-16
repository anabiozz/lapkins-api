package models

// Product ..
type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"decription"`
	Price       string `json:"price"`
	Size        string `json:"size"`
}

// Variant ..
type Variant struct {
	VariantID     int      `json:"variant_id"`
	ProductID     int      `json:"product_id"`
	Name          string   `json:"name"`
	Description   string   `json:"decription"`
	PriceOverride int      `json:"price_override"`
	Attributes    []string `json:"attributes"`
	Sizes         []string `json:"sizes"`
	Size          string   `json:"size"`
	Images        []string `json:"images"`
	Quantity      int      `json:"quantity"`
}
