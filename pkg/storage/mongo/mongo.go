package mongo

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Logger log.Logger
}

type Storage struct {
	db     *mongo.Database
	logger log.Logger
	mu     sync.Mutex
}

// GetProducts ..
func (m *Storage) GetProducts(ctx context.Context, category string) ([]*model.CatalogProduct, error) {
	var products []*model.CatalogProduct
	cur, err := m.db.Collection("products").Find(ctx, bson.M{"category": category})
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
func (p *Storage) GetVariation(ctx context.Context, sku string) (*model.DescriptionProduct, error) {
	descProduct := &model.DescriptionProduct{}
	var sizes []*model.Size
	product := &model.Product{}

	err := p.db.Collection("products").FindOne(ctx, bson.D{{"variations.sku", sku}}).Decode(product)
	if err != nil {
		return nil, err
	}

	for _, variation := range product.Variations {
		size := &model.Size{}
		if variation.SKU == sku {
			descProduct.Variation = variation
		}
		size.Key = variation.SKU
		size.Value = fmt.Sprintf("%dx%d", variation.Dimensions.Width, variation.Dimensions.Height)
		sizes = append(sizes, size)
	}

	descProduct.Attributes = product.Attributes
	descProduct.Name = product.Name
	descProduct.Brand = product.Brand
	descProduct.Description = product.Description
	descProduct.Kind = product.Kind
	descProduct.Season = product.Season
	descProduct.Sizes = sizes

	return descProduct, nil
}

// GetCategories ..
func (p *Storage) GetCategories(categoryURL string) (categories []model.Category, err error) {
	return nil, nil
}

// AddProduct ..
func (s *Storage) AddProduct(ctx context.Context, sku string, user *model.CartUser) (*model.CartUser, error) {

	if sku == "" {
		return nil, errors.New("sku should not be empty")
	}

	var err error
	cart := &model.Cart{}
	cartProduct := &model.CartProduct{}
	var objID primitive.ObjectID
	var switchID bool

	if user.ID == "" && user.TmpID == "" {
		s.mu.Lock()
		objID = primitive.NewObjectIDFromTimestamp(time.Now())
		user.TmpID = objID.Hex()
		s.mu.Unlock()
	} else if user.ID != "" && user.TmpID != "" {
		switchID = true
		objID, err = primitive.ObjectIDFromHex(user.TmpID)
		if err != nil {
			return nil, err
		}
	} else if user.ID != "" {
		objID, err = primitive.ObjectIDFromHex(user.ID)
		if err != nil {
			return nil, err
		}
	} else if user.TmpID != "" {
		objID, err = primitive.ObjectIDFromHex(user.TmpID)
		if err != nil {
			return nil, err
		}
	}

	product := &model.Product{}
	variation := &model.Variation{}

	err = s.db.Collection("products").FindOne(ctx, bson.D{{"variations.sku", sku}}).Decode(product)
	if err != nil {
		return nil, err
	}

	for _, v := range product.Variations {
		if v.SKU == sku {
			variation = v
		}
	}

	err = s.db.Collection("carts").FindOne(ctx, bson.D{{"_id", objID}, {"status", "active"}, {"products.sku", sku}}).Decode(cart)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {

			cartProduct.Name = product.Name
			cartProduct.Price = variation.Pricing.Retail
			cartProduct.Quantity = 1
			cartProduct.SKU = sku
			cartProduct.Size = fmt.Sprintf("%dx%d", variation.Dimensions.Width, variation.Dimensions.Height)
			cartProduct.UpdatedAt = time.Now()
			cartProduct.CreatedAt = time.Now()

			filter := bson.D{{"_id", objID}, {"status", "active"}}
			update := bson.D{
				{
					"$set",
					bson.D{
						{"updated_at", time.Now()},
						{"created_at", time.Now()},
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
			_, err := s.db.Collection("carts").UpdateOne(ctx, filter, update, opts)
			if err != nil {
				return nil, err
			}
		}
	} else {

		if switchID {
			newID, err := primitive.ObjectIDFromHex(user.ID)
			if err != nil {
				return nil, err
			}
			cart.ID = newID
			_, err = s.db.Collection("carts").InsertOne(ctx, cart)
			if err != nil {
				return nil, err
			}
			_, err = s.db.Collection("carts").DeleteOne(ctx, bson.D{{"_id", objID}, {"status", "active"}})
			if err != nil {
				return nil, err
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
			_, err := s.db.Collection("carts").UpdateOne(ctx, filter, update, opts)
			if err != nil {
				return nil, err
			}
		}
	}

	return user, nil
}

func (s Storage) GetHeaderCartInfo(ctx context.Context, userID string) (*model.HeaderCartInfo, error) {
	info := &model.HeaderCartInfo{}
	cart := &model.Cart{}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	err = s.db.Collection("carts").FindOne(ctx, bson.D{{"_id", objID}, {"status", "active"}}).Decode(cart)
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
	_, err = s.db.Collection("carts").UpdateOne(ctx, filter, update, opts)
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
	_, err = s.db.Collection("carts").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

// RemoveProduct ..
func (p *Storage) RemoveProduct(variationID int, cartSession string, sizeOptionID int) (err error) {
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
	err = s.db.Collection("carts").FindOne(ctx, bson.D{{"_id", objID}, {"status", "active"}}).Decode(cart)
	if err != nil {
		return nil, err
	}
	for _, product := range cart.Products {
		cartProducts = append(cartProducts, product)
	}
	return cartProducts, nil
}

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

func (m Storage) GetUserBySubject(ctx context.Context, email string, phone int64) (*model.User, error) {
	user := &model.User{}
	filter := bson.D{{"$or", bson.A{bson.D{{"email", email}}, bson.D{{"phone", phone}}}}}
	err := m.db.Collection("users").FindOne(ctx, filter).Decode(user)
	if err != nil {
		return user, errors.New("invalid subject")
	}
	return user, nil
}

func (m Storage) GetUsers(ctx context.Context) ([]*model.User, error) {
	filter := bson.D{}
	var users []*model.User
	cur, err := m.db.Collection("users").Find(ctx, filter)
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

func New(cfg Config) (*Storage, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	m := &Storage{
		db:     client.Database("lapkins"),
		logger: cfg.Logger,
	}
	level.Info(cfg.Logger).Log("msg", "mongo was up")
	return m, nil
}
