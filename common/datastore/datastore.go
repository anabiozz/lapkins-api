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
