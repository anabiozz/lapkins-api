package model

import (
	"time"
)

type CatalogProduct struct {
	LName     string  `json:"lname"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Thumbnail string  `json:"thumbnail"`
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
	ID           string         `bson:"_id" json:"id,omitempty"`
	Name         string         `json:"name"`
	LName        string         `json:"lname"`
	Descriptions []*Description `json:"desc"`
	Brand        *Brand         `json:"brand"`
	Variations   []*Variation   `json:"variations"`
	CreatedOn    time.Time      `json:"createdOn"`
	ModifiedOn   time.Time      `json:"modifiedOn"`
}

type Brand struct {
	Country *Country `json:"country"`
	Img     *Img     `json:"img"`
	Name    string   `json:"name"`
}

type Description struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

type Img struct {
	Src   string `json:"src"`
	Title string `json:"title"`
}

type Country struct {
	Name string `json:"name"`
}

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Size struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Variation ..
type Variation struct {
	Shipping   *Shipping    `json:"shipping"`
	Pricing    *Pricing     `json:"pricing"`
	Assets     *Assets      `json:"assets"`
	Attributes []*Attribute `json:"attributes"`
}

type Assets struct {
	Thumbnail *Thumbnail `json:"thumbnail"`
	Imgs      []*Img     `json:"imgs"`
}

type Thumbnail struct {
	Src string `json:"src"`
}

type Weight struct {
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

type Dimensions struct {
	Width  string `json:"width"`
	Height string `json:"height"`
	Length string `json:"length"`
	Unit   string `json:"unit"`
}

type Pricing struct {
	Price float64 `json:"price"`
	Sale  *Sale   `json:"sale"`
}

type Sale struct {
	SalePrice   float64 `json:"sale_price"`
	SaleEndDate string  `json:"sale_end_date"`
}
