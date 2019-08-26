package models

import "encoding/json"

// Product ..
type Product struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	Categories   json.RawMessage `json:"categories"`
	Currency     string          `json:"currency"`
	Description  string          `json:"decription"`
	Price        string          `json:"price"`
	ProductsType string          `json:"products_type"`
	IsAvailable  bool            `json:"is_available"`
	Ext          string          `json:"ext"`
}
