package mongo

import (
	"context"
	"fmt"

	"github.com/anabiozz/core/lapkins/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

// GetCategories ..
func (s *Storage) GetCategories(ctx context.Context) ([]*model.Category, error) {
	cursor, err := s.db.Collection("categories").Find(ctx, bson.D{{"parents", nil}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var categories []*model.Category
	for cursor.Next(ctx) {

		var category *model.Category
		if err = cursor.Decode(&category); err != nil {
			return nil, err
		}

		cursor, err := s.db.Collection("categories").Find(ctx, bson.D{{"parents", category.ID}})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var subcategory *model.Subcategory
			if err = cursor.Decode(&subcategory); err != nil {
				return nil, err
			}
			category.Ancestors = append(category.Ancestors, subcategory)
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
