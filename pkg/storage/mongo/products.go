package mongo

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetCatalog ..
func (s *Storage) GetCatalog(ctx context.Context, category string) ([]*model.CatalogProduct, error) {
	var products []*model.CatalogProduct

	findOptions := options.Find()
	findOptions.SetSkip(0)
	findOptions.SetLimit(10)

	cur, err := s.db.Collection("products").Find(ctx, bson.M{"category": category}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {

		product := &model.Product{}
		catalogProduct := &model.CatalogProduct{}

		err := cur.Decode(&product)
		if err != nil {
			return nil, err
		}

		sort.SliceStable(product.Variations, func(i, j int) bool {
			return product.Variations[i].Pricing.Retail < product.Variations[j].Pricing.Retail
		})

		catalogProduct.Name = product.Name
		catalogProduct.Price = product.Variations[0].Pricing.Retail
		catalogProduct.SKU = product.Variations[0].SKU

		products = append(products, catalogProduct)

	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

// GetVariation ..
func (s *Storage) GetProduct(ctx context.Context, sku string) (*model.Product, error) {
	product := &model.Product{}
	err := s.db.Collection("products").FindOne(ctx, bson.D{{"variations.sku", sku}}).Decode(product)
	if err != nil {
		return nil, err
	}
	for _, variation := range product.Variations {
		if variation.SKU == sku {
			product.Variations = nil
			product.Variations = make([]*model.Variation, 0, 1)
			product.Variations = append(product.Variations, variation)
		}

		size := &model.Size{}
		size.Key = variation.SKU
		size.Value = fmt.Sprintf("%dx%d", variation.Dimensions.Width, variation.Dimensions.Height)
		product.Sizes = append(product.Sizes, size)
	}
	return product, nil
}

func (s *Storage) GetProductsByCategory(ctx context.Context, category string) ([]*model.SKUProduct, error) {
	var products []*model.SKUProduct

	findOptions := options.Find()
	findOptions.SetSkip(0)
	findOptions.SetLimit(10)

	var filter bson.M
	if category != "" {
		filter = bson.M{"category": category}
	}
	cur, err := s.db.Collection("products").Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		product := &model.Product{}
		err := cur.Decode(&product)
		if err != nil {
			return nil, err
		}
		for _, variation := range product.Variations {
			skuproduct := &model.SKUProduct{}
			skuproduct.Season = product.Season
			skuproduct.Description = product.Description
			skuproduct.Brand = product.Brand
			skuproduct.Attributes = product.Attributes
			skuproduct.Kind = product.Kind
			skuproduct.ModifiedOn = time.Now()
			skuproduct.Category = product.Category
			skuproduct.Name = product.Name + " " + fmt.Sprintf("%dx%d", variation.Dimensions.Width, variation.Dimensions.Height)
			skuproduct.Variation = variation
			products = append(products, skuproduct)
		}
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (s *Storage) AddAttribute(ctx context.Context, sku string, attribute model.Attribute) error {
	fmt.Println("AddAttribute")
	filter := bson.D{{"status", "active"}, {"products.sku", sku}}
	update := bson.D{
		{
			"$set",
			bson.D{
				{"products.$.updated_at", time.Now()},
			},
		},
		{
			"$inc",
			bson.D{
				{"products.$.quantity", 1},
			},
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err = s.db.Collection("carts").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) RemoveAttribute(ctx context.Context, sku string, attribute string) error {
	fmt.Println("RemoveAttribute")
	return nil
}

func (s *Storage) AddCategory(ctx context.Context, sku string) error {
	fmt.Println("AddCategory")
	return nil
}

func (s *Storage) RemoveCategory(ctx context.Context, sku string) error {
	fmt.Println("RemoveCategory")
	return nil
}
