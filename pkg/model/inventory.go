package model

import "time"

type Inventory struct {
	ID           string `bson:"_id" json:"id,omitempty"`
	Quantity     int    `json:"quantity"`
	Reservations []*Reservation
	CreatedOn    time.Time `json:"createdOn"`
	ModifiedOn   time.Time `json:"modifiedOn"`
}

type Reservation struct {
	ID         string    `bson:"_id" json:"id,omitempty"`
	Quantity   int       `json:"quantity"`
	CreatedOn  time.Time `json:"createdOn"`
	ModifiedOn time.Time `json:"modifiedOn"`
}
