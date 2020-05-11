package mongo

import (
	"context"
	"github.com/anabiozz/lapkins-api/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetCatalog ..
func (s *Storage) GetCatalog(ctx context.Context, department string, category string) ([]*model.CatalogProduct, error) {
	var products []*model.CatalogProduct

	productsCur, err := s.db.Collection("products").Find(ctx, bson.D{{"dep", department}}, nil)
	if err != nil {
		return nil, err
	}
	defer productsCur.Close(ctx)

	for productsCur.Next(ctx) {

		product := &model.Product{}
		catalogProduct := &model.CatalogProduct{}

		err := productsCur.Decode(&product)
		if err != nil {
			return nil, err
		}

		var variation model.Variation
		findOptions := options.FindOne()
		findOptions.SetSkip(0)
		findOptions.SetSort(bson.M{"pricing.price": 1})
		err = s.db.Collection("variations").FindOne(ctx, bson.D{
			{"category", primitive.Regex{Pattern: category, Options: "i"}},
			{"productId", product.ID}}, findOptions).Decode(&variation)
		if err != nil {
			return nil, err
		}
		catalogProduct.ID = product.ID
		catalogProduct.Name = product.Name
		catalogProduct.LName = product.LName
		catalogProduct.Price = variation.Pricing.Price
		catalogProduct.Thumbnail = variation.Assets.Thumbnail.Src

		products = append(products, catalogProduct)
	}
	if err := productsCur.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

// GetVariation ..
func (s *Storage) GetProduct(ctx context.Context, sku string, attr string) (*model.VariationProduct, error) {
	product := &model.Product{}
	variationProduct := &model.VariationProduct{}

	err := s.db.Collection("products").FindOne(ctx, bson.D{{"_id", sku}}).Decode(product)
	if err != nil {
		return nil, err
	}

	variationProduct.Descriptions = product.Descriptions
	variationProduct.Name = product.Name
	variationProduct.Brand = product.Brand

	var variation model.Variation

	var filter bson.D
	if attr != "" {
		filter = bson.D{{"productId", product.ID}, {"attributes.value", attr}}
	} else {
		filter = bson.D{{"productId", product.ID}}
	}

	err = s.db.Collection("variations").FindOne(ctx, filter).Decode(&variation)
	if err != nil {
		return nil, err
	}

	variationProduct.Variation = &variation

	for _, attr := range product.Attrs {
		facetsCur, err := s.db.Collection("facets").Find(ctx, bson.D{{"name", attr.Name}}, nil)
		if err != nil {
			return nil, err
		}
		defer facetsCur.Close(ctx)

		var variationType model.VariationType

		for facetsCur.Next(ctx) {
			var facet model.Facet
			err := facetsCur.Decode(&facet)
			if err != nil {
				return nil, err
			}
			variationType.Name = facet.Name
			variationType.Display = facet.Display
			variationType.Attrs = append(variationType.Attrs, facet.Value)
		}

		variationProduct.VariationTypes = append(variationProduct.VariationTypes, &variationType)

		if err := facetsCur.Err(); err != nil {
			return nil, err
		}
	}
	return variationProduct, nil
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
		//for _, variation := range product.Variations {
		//	skuproduct := &model.SKUProduct{}
		//	skuproduct.Season = product.Season
		//	skuproduct.Description = product.Descriptions
		//	skuproduct.Brand = product.Brand
		//	skuproduct.Attributes = product.Attributes
		//	skuproduct.Kind = product.Kind
		//	skuproduct.ModifiedOn = time.Now()
		//	skuproduct.Category = product.Category
		//	skuproduct.Name = product.Name + " " + fmt.Sprintf("%dx%d", variation.Dimensions.Width, variation.Dimensions.Height)
		//	skuproduct.Variation = variation
		//	products = append(products, skuproduct)
		//}
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return products, nil
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
	_, err := s.db.Collection("carts").UpdateOne(ctx, filter, update, nil)
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
	_, err := s.db.Collection("carts").UpdateOne(ctx, filter, update, nil)
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
//	_, err := s.db.Collection("carts").UpdateOne(ctx, filter, update, nil)
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
