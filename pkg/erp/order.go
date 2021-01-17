package erp

import "time"

type Order struct {
	ID         string `bson:"_id" json:"id,omitempty"`
	TotalPrice int    `json:"total_price"`
	Shipping   *Shipping
	Pricing    *Pricing
	CreatedOn  time.Time `json:"createdOn"`
	ModifiedOn time.Time `json:"modifiedOn"`
}

type Shipping struct {
	Dimensions *Dimensions `json:"dimensions"`
	Weight     *Weight     `json:"weight"`
	Address    string      `json:"address"`
	CreatedOn  time.Time   `json:"createdOn"`
	ModifiedOn time.Time   `json:"modifiedOn"`
}

type Payment struct {
	ID            string    `bson:"_id" json:"id,omitempty"`
	Method        string    `json:"method"`
	TransactionID int       `json:"transaction_id"`
	CreatedOn     time.Time `json:"createdOn"`
	ModifiedOn    time.Time `json:"modifiedOn"`
}
