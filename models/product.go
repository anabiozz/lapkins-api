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
	Article     string `json:"aricle"`
	Price       string `json:"price"`
}

// Variant ..
type Variant struct {
	VariantID   int      `json:"variant_id"`
	ProductID   int      `json:"product_id"`
	Name        string   `json:"name"`
	Description string   `json:"decription"`
	Brand       string   `json:"brand"`
	Subject     string   `json:"subject"`
	Season      string   `json:"season"`
	Kind        string   `json:"kind"`
	Images      []string `json:"images"`
	Attributes  []string `json:"attributes"`
	Price       string   `json:"price"`
}
