package models

import "encoding/json"

// Product ..
type Product struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"decription"`
	Price              string `json:"price"`
	ProductVariantName string `json:"product_variant_name"`
}

// ProductVariant ..
type ProductVariant struct {
	ID            int             `json:"id"`
	ProductID     int             `json:"product_id"`
	Name          string          `json:"name"`
	Description   string          `json:"decription"`
	PriceOverride string          `json:"price_override"`
	Attributes    json.RawMessage `json:"attributes"`
	Sizes         []string        `json:"sizes"`
}
