package mongo

import (
	"context"
	"errors"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s Storage) GetUserBySubject(ctx context.Context, email string, phone int64) (*model.User, error) {
	user := &model.User{}
	filter := bson.D{{"$or", bson.A{bson.D{{"email", email}}, bson.D{{"phone", phone}}}}}
	err := s.db.Collection("users").FindOne(ctx, filter).Decode(user)
	if err != nil {
		return user, errors.New("invalid subject")
	}
	return user, nil
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
