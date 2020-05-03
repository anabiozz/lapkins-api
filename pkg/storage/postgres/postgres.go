package postgres

import (
	"database/sql"
	"fmt"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/kelseyhightower/envconfig"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// Configs for database
type Config struct {
	DBInfo string `envconfig:"lapkinenv" required:"true"`
	Logger *log.Logger
}

// Storage ..
type Storage struct {
	*sql.DB
	Logger *log.Logger
}

// NewStorage ..
func NewStorage(cfg *Config) (*Storage, error) {
	if err := envconfig.Process("", cfg); err != nil {
		cfg.Logger.Fatal(err)
	}

	connection, err := sql.Open("postgres", cfg.DBInfo)
	if err != nil {
		return nil, err
	}
	return &Storage{
		DB:     connection,
		Logger: cfg.Logger,
	}, nil
}

// GetCatalog ..
func (p *Storage) GetCatalog(categoryURL string) (products []model.Product, err error) {
	query := fmt.Sprintf(`SELECT * FROM products.get_products('%s');`, categoryURL)
	rows, err := p.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	product := model.Product{}
	for rows.Next() {

		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Brand,
			&product.Subject,
			&product.Season,
			&product.Kind,
			&product.PhotoCount,
			&product.Article,
			&product.Price,
			&product.CategoryDescription,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

// GetVariation ..
func (p *Storage) GetVariation(variationID, sizeOptionID string) (*model.Variation, error) {
	query := fmt.Sprintf(`SELECT * FROM products.get_variation(%s, %s);`, variationID, sizeOptionID)

	variation := &model.Variation{}

	err := p.QueryRow(query).Scan(
		&variation.ID,
		&variation.ProductID,
		&variation.Name,
		&variation.Description,
		&variation.Brand,
		&variation.Subject,
		&variation.Season,
		&variation.Kind,
		pq.Array(&variation.Images),
		pq.Array(&variation.Attributes),
		&variation.Price,
		pq.Array(&variation.Sizes),
		&variation.SizeOptionID,
	)

	if err != nil {
		return nil, err
	}

	return variation, nil
}

// CloseDB ..
func (p *Storage) CloseDB() {
	p.DB.Close()
}

// GetCategories ..
func (p *Storage) GetCategories(categoryURL string) (categories []model.Category, err error) {
	query := fmt.Sprintf(`SELECT * FROM products.get_categories('%s');`, categoryURL)
	rows, err := p.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	category := model.Category{}
	for rows.Next() {

		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.URL,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// AddProduct ..
func (p *Storage) AddProduct(variationID int, cartSession string, sizeOptionID int) (err error) {
	query := fmt.Sprintf(`SELECT * FROM carts.add_product(%d, '%s', %d);`, variationID, cartSession, sizeOptionID)
	_, err = p.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// IncreaseProductQuantity ..
func (p *Storage) IncreaseProductQuantity(variationID int, cartSession string, sizeOptionID int) (err error) {
	query := fmt.Sprintf(`SELECT * FROM carts.add_product(%d, '%s', %d);`, variationID, cartSession, sizeOptionID)
	_, err = p.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// DecreaseProductQuantity ..
func (p *Storage) DecreaseProductQuantity(variationID int, cartSession string, sizeOptionID int) (err error) {
	query := fmt.Sprintf(`SELECT * FROM carts.decrease_product_quantity(%d, '%s', %d);`, variationID, cartSession, sizeOptionID)
	_, err = p.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// RemoveProduct ..
func (p *Storage) RemoveProduct(variationID int, cartSession string, sizeOptionID int) (err error) {
	query := fmt.Sprintf(`SELECT * FROM carts.remove_product(%d, '%s', %d);`, variationID, cartSession, sizeOptionID)
	_, err = p.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// LoadCart ..
func (p *Storage) LoadCart(cartSession string) (cartItems []model.CartItemResponse, err error) {
	query := fmt.Sprintf(`SELECT * FROM carts.load_cart('%s');`, cartSession)

	rows, err := p.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cartItem := model.CartItemResponse{}
	for rows.Next() {

		err = rows.Scan(
			&cartItem.ID,
			&cartItem.Name,
			&cartItem.Brand,
			&cartItem.Price,
			&cartItem.PricePerItem,
			&cartItem.Size,
			&cartItem.Quantity,
			&cartItem.SizeOptionID,
		)
		if err != nil {
			return nil, err
		}

		cartItems = append(cartItems, cartItem)
	}

	if err != nil {
		return cartItems, err
	}

	return cartItems, nil
}
