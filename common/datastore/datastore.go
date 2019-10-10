package datastore

import (
	"errors"

	"github.com/anabiozz/lapkin-project/lapkin-api/models"
)

// Datastore ...
type Datastore interface {
	GetProducts(productsID string) (products []models.Product, err error)
	GetVariant(productVariantID, size string) (product *models.Variant, err error)
	GetCategories(categoryID string) (categories models.Categories, err error)

	AddProduct(variant []byte) (cartSession string, err error)
	ChangeQuantity(variantID string, cartSession string, newQuantety string) (err error)
	RemoveProduct(cartSession string, variant *models.Variant) (err error)
	GetCart(cartSession string) (cartItems []*models.Variant, err error)

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
