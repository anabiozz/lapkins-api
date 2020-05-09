package model

import "time"

type CatalogProduct struct {
	SKU   string `json:"sku"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type SKUProduct struct {
	Category    string       `json:"category"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Brand       string       `json:"brand"`
	Season      string       `json:"season"`
	Kind        string       `json:"kind"`
	Attributes  []*Attribute `json:"attributes"`
	Sizes       []*Size      `json:"sizes"`
	Variation   *Variation   `json:"variation"`
	CreatedOn   time.Time    `json:"createdOn"`
	ModifiedOn  time.Time    `json:"modifiedOn"`
}

// Product ..
type Product struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Brand       *Brand       `json:"brand"`
	Attributes  []*Attribute `json:"attributes"`
	Variations  []*Variation `json:"variations"`
	CreatedOn   time.Time    `json:"createdOn"`
	ModifiedOn  time.Time    `json:"modifiedOn"`
}

type Brand struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Attribute struct {
	Name string `json:"name"`
	Ext  *Ext   `json:"ext"`
}

type Ext struct {
	Name   string `json:"name"`
	IsHide bool   `json:"is_hide"`
}

type Size struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Variation ..
type Variation struct {
	Name   string `json:"name"`
	Exts   []*Ext `json:"exts"`
	IsHide bool   `json:"is_hide"`
}

type Weight struct {
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

type Dimensions struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Depth  int    `json:"depth"`
	Unit   string `json:"unit"`
}

type Pricing struct {
	List       int `json:"list"`
	Retail     int `json:"retail"`
	Savings    int `json:"savings"`
	PctSavings int `json:"pct_savings"`
}
