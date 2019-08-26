package models

// Currency ..
type Currency struct {
	RU int8
	EU int8
	US int8
}

const (
	// RU ..
	RU = iota
	// EU ..
	EU
	// US ..
	US
)
