package datastore

import (
	"errors"

	"github.com/anabiozz/lapkin-project/lapkin/backend/models"
)

// Datastore ...
type Datastore interface {
	GetProducts(productsID string, paths models.Paths) (products []models.Product, err error)
	GetProductByID(productID string) (product *models.Product, err error)
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
