package mongo

import (
	"context"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

// GetCategories ..
func (s *Storage) GetCategory(ctx context.Context, category string) ([]*model.Category, error) {
	cursor, err := s.db.Collection("categories").Find(ctx, bson.D{{"parent", category}})
	if err != nil {
		return nil, err
	}
	var categories []*model.Category
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var category *model.Category
		if err = cursor.Decode(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
