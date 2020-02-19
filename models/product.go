package models

// Product ..
type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"decription"`
	Brand       string `json:"brand"`
	Subject     string `json:"subject"`
	Season      string `json:"season"`
	Kind        string `json:"kind"`
	PhotoCount  string `json:"photo_count"`
	Article     string `json:"article"`
	Price       string `json:"price"`
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
