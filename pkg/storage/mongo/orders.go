package mongo

import (
	"context"

	"github.com/anabiozz/lapkins-api/pkg/model"
)

func (s *Storage) AddOrder(ctx context.Context, order *model.Order) error {
	_, err := s.db.Collection("orders").InsertOne(ctx, order, nil)
	if err != nil {
		return err
	}
	return nil
}
