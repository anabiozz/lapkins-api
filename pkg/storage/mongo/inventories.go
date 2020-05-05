package mongo

import (
	"context"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) AddReservation(ctx context.Context, sku string, category *model.Category) error {
	filter := bson.D{{"status", "active"}, {"products.sku", sku}}
	update := bson.D{
		{
			"$push",
			bson.D{
				{"products", category},
			},
		},
	}
	_, err := s.db.Collection("reservation").UpdateOne(ctx, filter, update, nil)
	if err != nil {
		return err
	}
	return nil
}
