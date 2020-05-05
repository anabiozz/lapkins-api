package model

import "time"

// Category ..
type Category struct {
	ID          string    `bson:"_id" json:"id,omitempty"`
	Ancestors   []string  `json:"ancestors"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Parent      string    `json:"parent"`
	CreatedOn   time.Time `json:"createdOn"`
	ModifiedOn  time.Time `json:"modifiedOn"`
}
