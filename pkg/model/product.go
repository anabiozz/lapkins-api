package model

import "time"

type CatalogProduct struct {
	SKU   string `json:"sku"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type DescriptionProduct struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Brand       string      `json:"brand"`
	Season      string      `json:"season"`
	Kind        string      `json:"kind"`
	Sizes       []*Size     `json:"sizes"`
	Attributes  []Attribute `json:"attributes"`
	Variation   *Variation  `json:"variation"`
}

// Product ..
type Product struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Brand       string       `json:"brand"`
	Season      string       `json:"season"`
	Kind        string       `json:"kind"`
	Attributes  []Attribute  `json:"attributes"`
	Sizes       []Size       `json:"sizes"`
	Variations  []*Variation `json:"variations"`
	CreatedOn   time.Time    `json:"createdOn"`
	ModifiedOn  time.Time    `json:"modifiedOn"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Size struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Variation ..
type Variation struct {
	SKU        string     `json:"sku"`
	Weight     Weight     `json:"weight"`
	Dimensions Dimensions `json:"dimensions"`
	Pricing    Pricing    `json:"pricing"`
	Photos     []string   `json:"photos"`
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

////JSONRaw ...
//type JSONRaw json.RawMessage
//
////Value ...
//func (j JSONRaw) Value() (driver.Value, error) {
//	byteArr := []byte(j)
//
//	return driver.Value(byteArr), nil
//}
//
////Scan ...
//func (j *JSONRaw) Scan(src interface{}) error {
//	asBytes, ok := src.([]byte)
//	if !ok {
//		return error(errors.New("scan source was not []bytes"))
//	}
//	err := json.Unmarshal(asBytes, &j)
//	if err != nil {
//		return error(errors.New("scan could not unmarshal to []string"))
//	}
//
//	return nil
//}
//
////MarshalJSON ...
//func (j *JSONRaw) MarshalJSON() ([]byte, error) {
//	return *j, nil
//}
//
////UnmarshalJSON ...
//func (j *JSONRaw) UnmarshalJSON(data []byte) error {
//	if j == nil {
//		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
//	}
//	*j = append((*j)[0:0], data...)
//	return nil
//}
//
//// Item ..
//type Item struct {
//	Key   int    `json:"key"`
//	Value string `json:"value"`
//}
