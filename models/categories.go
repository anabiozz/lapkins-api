package models

import "encoding/json"

// Category ..
type Category struct {
	Display string `json:"display"`
	URL     string `json:"url"`
}

// Categories ..
type Categories struct {
	Name     string          `json:"name"`
	Category json.RawMessage `json:"categories"`
}
