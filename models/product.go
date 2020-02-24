package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Product ..
type Product struct {
	ID                  int    `json:"variation_id"`
	Name                string `json:"name"`
	Description         string `json:"decription"`
	Brand               string `json:"brand"`
	Subject             string `json:"subject"`
	Season              string `json:"season"`
	Kind                string `json:"kind"`
	PhotoCount          string `json:"photo_count"`
	Article             string `json:"article"`
	Price               string `json:"price"`
	CategiryDescription string `json:"category_descrption"`
}

//JSONRaw ...
type JSONRaw json.RawMessage

//Value ...
func (j JSONRaw) Value() (driver.Value, error) {
	byteArr := []byte(j)

	return driver.Value(byteArr), nil
}

//Scan ...
func (j *JSONRaw) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes"))
	}
	err := json.Unmarshal(asBytes, &j)
	if err != nil {
		return error(errors.New("Scan could not unmarshal to []string"))
	}

	return nil
}

//MarshalJSON ...
func (j *JSONRaw) MarshalJSON() ([]byte, error) {
	return *j, nil
}

//UnmarshalJSON ...
func (j *JSONRaw) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

// Item ..
type Item struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}

// Variation ..
type Variation struct {
	ID           int       `json:"variation_id"`
	ProductID    int       `json:"product_id"`
	Name         string    `json:"name"`
	Description  string    `json:"decription"`
	Brand        string    `json:"brand"`
	Subject      string    `json:"subject"`
	Season       string    `json:"season"`
	Kind         string    `json:"kind"`
	Images       []string  `json:"images"`
	Attributes   []JSONRaw `json:"attributes"`
	Sizes        []JSONRaw `json:"sizes"`
	Price        string    `json:"price"`
	SizeOptionID int       `json:"size_option_id"`
}
