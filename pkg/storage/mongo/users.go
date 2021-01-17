package mongo

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/anabiozz/core/lapkins/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s Storage) RegisterUser(ctx context.Context, user *model.User) (string, error) {
	err := s.db.Collection("users").FindOne(ctx, bson.D{{"email", user.Email}, {"phone", user.Phone}}).Decode(&model.User{})
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			result, err := s.db.Collection("users").InsertOne(ctx, user)
			if err != nil {
				return "", errors.New("error while creating users")
			}
			if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
				return oid.Hex(), nil
			}
			return result.InsertedID.(string), nil
		}
		return "", err
	}
	return "", errors.New("this users already exists")
}

func (s Storage) Login(ctx context.Context, email string, phone int64, tmpUserID string) (*model.User, error) {

	user := &model.User{}
	filter := bson.D{{"$or", bson.A{bson.D{{"email", email}}, bson.D{{"phone", phone}}}}}
	err := s.db.Collection("users").FindOne(ctx, filter).Decode(user)
	if err != nil {
		return user, errors.New("invalid subject")
	}

	if tmpUserID != "" {

		cart := &model.Cart{}
		tmpCart := &model.Cart{}
		tmpCartID, err := primitive.ObjectIDFromHex(tmpUserID)
		// ищу корзину с временной айдихой
		err = s.db.Collection("cart").FindOne(ctx, bson.D{{"_id", tmpCartID}, {"status", "active"}}).Decode(tmpCart)
		if err != nil {
			return nil, err
		}
		_, err = s.db.Collection("cart").DeleteOne(ctx, bson.D{{"_id", tmpCartID}, {"status", "active"}})
		if err != nil {
			return nil, err
		}

		// Ищу корзину с айдихой юзера
		err = s.db.Collection("cart").FindOne(ctx, bson.D{{"_id", user.ID}, {"status", "active"}}).Decode(cart)
		if err != nil {
			if err.Error() != "mongo: no documents in result" {
				return nil, err
			}
			tmpCart.ID = user.ID
			_, err = s.db.Collection("cart").InsertOne(ctx, tmpCart)
			if err != nil {
				return nil, err
			}
		} else {

			sort.SliceStable(cart.Products, func(i, j int) bool {
				return cart.Products[i].SKU < cart.Products[j].SKU
			})

			sort.SliceStable(tmpCart.Products, func(i, j int) bool {
				return tmpCart.Products[i].SKU < tmpCart.Products[j].SKU
			})

			var filter bson.D
			var update bson.D

			for i, tmpProduct := range tmpCart.Products {
				if contains(cart.Products, tmpProduct) {
					filter = bson.D{{"_id", cart.ID}, {"status", "active"}, {"products.sku", cart.Products[i].SKU}}
					update = bson.D{
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
				} else {
					tmpProduct.UpdatedAt = time.Now()
					tmpCart.UpdatedAt = time.Now()
					filter = bson.D{{"_id", cart.ID}, {"status", "active"}, {"products.sku", cart.Products[i].SKU}}
					update = bson.D{
						{
							"$push",
							bson.D{
								{"products", tmpProduct},
							},
						},
					}
				}

				opts := options.Update().SetUpsert(true)
				_, err = s.db.Collection("cart").UpdateOne(ctx, filter, update, opts)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return user, nil
}

func contains(s []*model.CartProduct, e *model.CartProduct) bool {
	for _, a := range s {
		if a.SKU == e.SKU {
			return true
		}
	}
	return false
}

func (s Storage) GetUsers(ctx context.Context) ([]*model.User, error) {
	filter := bson.D{}
	var users []*model.User
	cur, err := s.db.Collection("users").Find(ctx, filter)
	if err != nil {
		return nil, errors.New("invalid subject")
	}
	for cur.Next(context.TODO()) {
		user := &model.User{}
		err := cur.Decode(user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	cur.Close(context.TODO())
	return users, nil
}
