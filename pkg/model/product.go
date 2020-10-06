package model

import (
	"time"
)

type CatalogProduct struct {
	ID        string  `json:"id,omitempty"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Thumbnail string  `json:"thumbnail"`
}

type Product struct {
	ID          int    `json:"id,omitempty"`
	Category    string `json:"category"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Attributes  []struct {
		Name  string   `json:"name"`
		Value []string `json:"value,omitempty"`
	} `json:"attributes,omitempty"`
	Variations []*Variation `json:"variations"`
	Variation  *Variation   `json:"variation"`
	CreatedOn  time.Time    `json:"createdOn"`
	ModifiedOn time.Time    `json:"modifiedOn"`
}

type Variation struct {
	ID         int    `json:"id,omitempty"`
	ProductID  int    `json:"product_id"`
	SKU        int    `json:"sku"`
	Name       string `json:"name"`
	Display    string `json:"display"`
	Default    bool   `json:"default"`
	Attributes []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"attributes"`
	Price     string `json:"price"`
	Thumbnail string `json:"thumbnail"`
	Images    []struct {
		Src string `json:"src"`
	} `json:"images"`
	CreatedOn  time.Time `json:"createdOn"`
	ModifiedOn time.Time `json:"modifiedOn"`
}

type VariationType struct {
	Name    string   `bson:"name" json:"name"`
	Display string   `bson:"display" json:"display"`
	Attrs   []string `bson:"attrs" json:"attrs"`
}

type ProductAttr struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Facet struct {
	ID      string `bson:"_id" json:"id,omitempty"`
	Name    string `bson:"name" json:"name"`
	Display string `bson:"display" json:"display"`
	Value   string `bson:"value" json:"value"`
}

type Brand struct {
	Country *Country `json:"country"`
	Img     *Img     `json:"img"`
	Name    string   `json:"name"`
}

type LangValue struct {
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

type NameValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Size struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ProductVariation struct {
	ID         string   `bson:"_id" json:"id,omitempty"`
	Imgs       []*Img   `bson:"imgs" json:"imgs"`
	Attributes []string `bson:"attrs" json:"attrs"`
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
