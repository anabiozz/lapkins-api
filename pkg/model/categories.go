package model

// Category ..
type Category struct {
	ID          string `bson:"_id" json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
