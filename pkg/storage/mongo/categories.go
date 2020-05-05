package mongo

import (
	"context"
	"fmt"

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

func (s *Storage) AddCategory(ctx context.Context, sku string, category *model.Category) error {
	filter := bson.D{{"status", "active"}, {"products.sku", sku}}
	update := bson.D{
		{
			"$push",
			bson.D{
				{"products", category},
			},
		},
	}
	_, err := s.db.Collection("categories").UpdateOne(ctx, filter, update, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) RemoveCategory(ctx context.Context, sku string, category *model.Category) error {
	fmt.Println("RemoveCategory")
	return nil
}
