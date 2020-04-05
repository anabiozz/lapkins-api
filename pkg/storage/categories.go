package storage

// Category ..
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// Categories ..
// type Categories struct {
// 	Categories json.RawMessage `json:"categories"`
// }
