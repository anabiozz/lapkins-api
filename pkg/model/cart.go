package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CartItem ..
type Cart struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Status    string             `bson:"status" json:"status"`
	Products  []*CartProduct     `bson:"products" json:"products"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type CartProduct struct {
	Name      string    `json:"name"`
	SKU       string    `json:"sku"`
	Price     int       `json:"price"`
	Quantity  int       `json:"quantity"`
	Size      string    `json:"size"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type CartUser struct {
	ID    string
	TmpID string
}

type HeaderCartInfo struct {
	Price    int `json:"price"`
	Quantity int `json:"quantity"`
}
