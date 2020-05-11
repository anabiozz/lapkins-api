package model

import "time"

// Category ..
type Category struct {
	ID         string         `bson:"_id" json:"id,omitempty"`
	Name       []*LangValue   `json:"name"`
	Ancestors  []*Subcategory `json:"ancestors"`
	CreatedOn  time.Time      `json:"createdOn"`
	ModifiedOn time.Time      `json:"modifiedOn"`
}

type Subcategory struct {
	ID          string       `bson:"_id" json:"id,omitempty"`
	Name        []*LangValue `json:"name"`
	Description []*LangValue `json:"description"`
	Parents     []string     `json:"parents"`
	Facets      []string     `json:"facets"`
	CreatedOn   time.Time    `json:"createdOn"`
	ModifiedOn  time.Time    `json:"modifiedOn"`
}
