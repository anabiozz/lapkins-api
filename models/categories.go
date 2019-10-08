package models

import "encoding/json"

// Category ..
type Category struct {
	Display string `json:"display"`
	URL     string `json:"url"`
}

// Categories ..
type Categories struct {
	Categories json.RawMessage `json:"categories"`
}
