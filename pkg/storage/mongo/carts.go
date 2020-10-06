package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddProduct ..
func (s *Storage) AddProduct(ctx context.Context, sku string, userID string, isLoggedIn bool, isTmpUserIDSet bool) (bool, string, error) {
	if sku == "" {
		return false, "", errors.New("sku should not be empty")
	}

	var err error
	var setTmpUserIDCookie bool
	var cartID primitive.ObjectID
	var cartType string

	switch true {
	// Когда не зарегался и добавляю в корзину первый товар
	case userID == "" && !isLoggedIn:
		s.mu.Lock()
		cartID = primitive.NewObjectIDFromTimestamp(time.Now())
		userID = cartID.Hex()
		setTmpUserIDCookie = true
		cartType = "tmp"
		s.mu.Unlock()
	case userID != "":
		if isLoggedIn {
			cartType = "logged_in"
		}
		cartID, err = primitive.ObjectIDFromHex(userID)
		if err != nil {
			return false, "", err
		}
	}

	cart := &model.Cart{}
	product := &model.Product{}
	//variation := &model.Variation{}

	err = s.db.Collection("products").FindOne(ctx, bson.D{{"variations.sku", sku}}).Decode(product)
	if err != nil {
		return false, "", err
	}

	//for _, v := range product.Variations {
	//	if v.SKU == sku {
	//		variation = v
	//	}
	//}

	err = s.db.Collection("cart").FindOne(ctx, bson.D{{"_id", cartID}, {"status", "active"}, {"products.sku", sku}}).Decode(cart)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {

			if isLoggedIn && isTmpUserIDSet {

			}

			// Корзина не найдена
			cartProduct := &model.CartProduct{}
			cartProduct.Name = product.Name
			//cartProduct.Price = variation.Pricing.Retail
			cartProduct.Quantity = 1
			cartProduct.SKU = sku
			//cartProduct.Size = fmt.Sprintf("%dx%d", variation.Dimensions.Width, variation.Dimensions.Height)
			cartProduct.UpdatedAt = time.Now()
			cartProduct.CreatedAt = time.Now()

			filter := bson.D{{"_id", cartID}, {"status", "active"}}
			update := bson.D{
				{
					"$set",
					bson.D{
						{"updated_at", time.Now()},
						{"created_at", time.Now()},
						{"type", cartType},
					},
				},
				{
					"$push",
					bson.D{
						{"products", cartProduct},
					},
				},
			}
			opts := options.Update().SetUpsert(true)
			_, err := s.db.Collection("cart").UpdateOne(ctx, filter, update, opts)
			if err != nil {
				return false, "", err
			}
		}
	} else {

		filter := bson.D{{"_id", cart.ID}, {"status", "active"}, {"products.sku", sku}}
		update := bson.D{
			{
				"$set",
				bson.D{
					{"updated_at", time.Now()},
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
		_, err := s.db.Collection("cart").UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return false, "", err
		}
	}

	return setTmpUserIDCookie, userID, nil
}

func (s Storage) GetHeaderCartInfo(ctx context.Context, userID string) (*model.HeaderCartInfo, error) {
	info := &model.HeaderCartInfo{}
	cart := &model.Cart{}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	err = s.db.Collection("cart").FindOne(ctx, bson.D{{"_id", objID}, {"status", "active"}}).Decode(cart)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err
	}

	for _, product := range cart.Products {
		info.Price += product.Quantity * product.Price
		info.Quantity += product.Quantity
	}

	return info, nil
}

// IncreaseProductQuantity ..
func (s *Storage) IncreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", objID}, {"status", "active"}, {"products.sku", sku}}
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
	_, err = s.db.Collection("cart").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

// DecreaseProductQuantity ..
func (s *Storage) DecreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", objID}, {"status", "active"}, {"products.sku", sku}}
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
				{"products.$.quantity", -1},
			},
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err = s.db.Collection("cart").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

// RemoveProduct ..
func (s *Storage) RemoveProduct(ctx context.Context, userID string, sku string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", objID}, {"status", "active"}}
	update := bson.D{
		{
			"$pull",
			bson.D{
				{"products", bson.D{{"sku", sku}}},
			},
		},
		{
			"$set",
			bson.D{
				{"updated_at", time.Now()},
			},
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err = s.db.Collection("cart").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

// LoadCart ..
func (s *Storage) LoadCart(ctx context.Context, userID string) ([]*model.CartProduct, error) {
	cart := &model.Cart{}
	var cartProducts []*model.CartProduct

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	err = s.db.Collection("cart").FindOne(ctx, bson.D{{"_id", objID}, {"status", "active"}}).Decode(cart)
	if err != nil {
		return nil, err
	}
	for _, product := range cart.Products {
		cartProducts = append(cartProducts, product)
	}
	return cartProducts, nil
}
