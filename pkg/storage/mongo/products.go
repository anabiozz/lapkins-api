package mongo

import (
	"context"
	"strconv"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) GetProducts(ctx context.Context) ([]*model.Product, error) {
	var products []*model.Product
	productsCur, err := s.db.Collection("products").Find(ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	defer productsCur.Close(ctx)

	for productsCur.Next(ctx) {
		product := &model.Product{}
		err := productsCur.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := productsCur.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (s *Storage) UpdateProduct(ctx context.Context, product *model.Product) error {
	filter := bson.M{"id": bson.M{"$eq": product.ID}}
	update := bson.M{
		"$set": product,
	}
	_, err := s.db.Collection("products").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// GetVariation ..
func (s *Storage) GetProduct(ctx context.Context, sku string) (*model.Product, error) {
	var products []*model.Product
	resultProduct := &model.Product{}

	productsCur, err := s.db.Collection("products").Find(ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	defer productsCur.Close(ctx)

	for productsCur.Next(ctx) {
		product := &model.Product{}
		err := productsCur.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := productsCur.Err(); err != nil {
		return nil, err
	}

	skuInt, err := strconv.Atoi(sku)
	if err != nil {
		return nil, err
	}

	for _, product := range products {
		for _, variation := range product.Variations {
			if variation.SKU == skuInt {
				resultProduct.ID = product.ID
				resultProduct.Name = product.Name
				resultProduct.Category = product.Category
				resultProduct.Description = product.Description
				resultProduct.Attributes = product.Attributes
				resultProduct.Variation = variation
				resultProduct.Variations = product.Variations
			}
		}
	}

	return resultProduct, nil
}

func (s *Storage) AddAttribute(ctx context.Context, sku string, attribute *model.NameValue) error {
	filter := bson.D{{"status", "active"}, {"products.sku", sku}}
	update := bson.D{
		{
			"$push",
			bson.D{
				{"products", attribute},
			},
		},
	}
	_, err := s.db.Collection("cart").UpdateOne(ctx, filter, update, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) RemoveAttribute(ctx context.Context, sku string, attribute string) error {
	filter := bson.D{{"status", "active"}, {"products.sku", sku}}
	update := bson.D{
		{
			"$pull",
			bson.D{
				{"products", attribute},
			},
		},
	}
	_, err := s.db.Collection("cart").UpdateOne(ctx, filter, update, nil)
	if err != nil {
		return err
	}
	return nil
}

//func (s *Storage) AddCategory(ctx context.Context, sku string, category model.Category) error {
//	filter := bson.D{{"status", "active"}, {"products.sku", sku}}
//	update := bson.D{
//		{
//			"$push",
//			bson.D{
//				{"products", category},
//			},
//		},
//	}
//	_, err := s.db.Collection("cart").UpdateOne(ctx, filter, update, nil)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (s *Storage) RemoveCategory(ctx context.Context, sku string) error {
//	fmt.Println("RemoveCategory")
//	return nil
//}
