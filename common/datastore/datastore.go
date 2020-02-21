package datastore

import (
	"errors"

	"github.com/anabiozz/lapkins-api/models"
)

// Datastore ...
type Datastore interface {
	GetProducts(productsID string) (products []models.Product, err error)
	GetVariation(variationID, size string) (product *models.Variation, err error)
	GetCategories(categoryID string) (categories models.Categories, err error)

	CreateSession() (cartSession string, err error)
	AddProduct(variationID int, сartSession string) (err error)
	ChangeQuantity(variationID string, cartSession string, newQuantety string) (err error)
	RemoveProduct(cartSession string, variation *models.Variation) (err error)
	GetCart(cartSession string) (cartItems []*models.Variation, err error)

	CloseDB()
}

const (
	// POSTGRES ...
	POSTGRES = iota
)

// NewDatastore ...
func NewDatastore(datastoreType int) (Datastore, error) {
	switch datastoreType {
	case POSTGRES:
		return NewPostgresDatastore()
	}
	return nil, errors.New("unknown datastore type")
}
