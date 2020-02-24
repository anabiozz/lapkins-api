package datastore

import (
	"database/sql"
	"fmt"

	"github.com/anabiozz/lapkins-api/models"
	"github.com/anabiozz/logger"
	"github.com/kelseyhightower/envconfig"
	"github.com/lib/pq"
)

// Configs for database
type Configs struct {
	DBinfo string `envconfig:"lapkinenv" required:"true"`
}

// PostgresDatastore ..
type PostgresDatastore struct {
	*sql.DB
}

// NewPostgresDatastore ..
func NewPostgresDatastore() (*PostgresDatastore, error) {
	var db Configs
	if err := envconfig.Process("", &db); err != nil {
		logger.Fatal(err)
	}

	connection, err := sql.Open("postgres", db.DBinfo)
	if err != nil {
		return nil, err
	}
	return &PostgresDatastore{
		DB: connection,
	}, nil
}

// GetProducts ..
func (p *PostgresDatastore) GetProducts(categoryURL string) (products []models.Product, err error) {
	query := fmt.Sprintf(`SELECT * FROM products_v2.get_products('%s');`, categoryURL)
	rows, err := p.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	product := models.Product{}
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
			&product.CategiryDescription,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

// GetVariation ..
func (p *PostgresDatastore) GetVariation(variationID, sizeOptionID string) (*models.Variation, error) {
	query := fmt.Sprintf(`SELECT * FROM products_v2.get_variation(%s, %s);`, variationID, sizeOptionID)

	variation := &models.Variation{}

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
func (p *PostgresDatastore) CloseDB() {
	p.DB.Close()
}

// GetCategories ..
func (p *PostgresDatastore) GetCategories(categoryURL string) (categories []models.Category, err error) {
	query := fmt.Sprintf(`SELECT * FROM products_v2.get_categories('%s');`, categoryURL)
	rows, err := p.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	category := models.Category{}
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
func (p *PostgresDatastore) AddProduct(variationID int, cartSession string, sizeOptionID int) (err error) {
	query := fmt.Sprintf(`SELECT * FROM cart.add_product(%d, '%s', %d);`, variationID, cartSession, sizeOptionID)
	_, err = p.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// IncreaseProductQuantity ..
func (p *PostgresDatastore) IncreaseProductQuantity(variationID int, cartSession string, sizeOptionID int) (err error) {
	query := fmt.Sprintf(`SELECT * FROM cart.add_product(%d, '%s', %d);`, variationID, cartSession, sizeOptionID)
	_, err = p.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// DecreaseProductQuantity ..
func (p *PostgresDatastore) DecreaseProductQuantity(variationID int, cartSession string, sizeOptionID int) (err error) {
	query := fmt.Sprintf(`SELECT * FROM cart.decrease_product_quantity(%d, '%s', %d);`, variationID, cartSession, sizeOptionID)
	_, err = p.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// RemoveProduct ..
func (p *PostgresDatastore) RemoveProduct(variationID int, cartSession string, sizeOptionID int) (err error) {
	query := fmt.Sprintf(`SELECT * FROM cart.remove_product(%d, '%s', %d);`, variationID, cartSession, sizeOptionID)
	_, err = p.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// LoadCart ..
func (p *PostgresDatastore) LoadCart(cartSession string) (cartItems []models.CartItemResponse, err error) {
	query := fmt.Sprintf(`SELECT * FROM cart.load_cart('%s');`, cartSession)

	rows, err := p.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cartItem := models.CartItemResponse{}
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
