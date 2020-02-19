package datastore

import (
	"database/sql"
	"fmt"
	"strconv"

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
func (p *PostgresDatastore) GetProducts(subjectID string) (products []models.Product, err error) {
	id, err := strconv.Atoi(subjectID)
	query := fmt.Sprintf(`SELECT * FROM new_products.get_products(%d);`, id)
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
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

// GetVariant ..
func (p *PostgresDatastore) GetVariant(variantID, size string) (*models.Variant, error) {
	query := fmt.Sprintf(`SELECT * FROM new_products.get_variant(%s);`, variantID)

	variant := &models.Variant{}

	err := p.QueryRow(query).Scan(
		&variant.VariantID,
		&variant.ProductID,
		&variant.Name,
		&variant.Description,
		&variant.PriceOverride,
		pq.Array(&variant.Sizes),
		pq.Array(&variant.Images),
		pq.Array(&variant.Attributes),
	)

	if err != nil {
		return nil, err
	}

	return variant, nil
}

// CloseDB ..
func (p *PostgresDatastore) CloseDB() {
	p.DB.Close()
}

// GetCategories ..
func (p *PostgresDatastore) GetCategories(categoryID string) (models.Categories, error) {
	query := fmt.Sprintf(`SELECT * FROM products.get_categories(%s);`, categoryID)

	categories := models.Categories{}

	err := p.QueryRow(query).Scan(
		&categories.Categories,
	)

	if err != nil {
		return categories, err
	}

	return categories, nil
}

// CreateSession ..
func (p *PostgresDatastore) CreateSession() (cartSession string, err error) {
	query := fmt.Sprintf(`SELECT * FROM cart.create_session();`)
	err = p.QueryRow(query).Scan(
		&cartSession,
	)
	if err != nil {
		return "", err
	}
	return cartSession, nil
}

// AddProduct ..
func (p *PostgresDatastore) AddProduct(variantID, сartSession, customerID int) (cartSession string, err error) {
	query := fmt.Sprintf(`SELECT * FROM cart.add_product(%d, %d, %d);`, variantID, сartSession, customerID)
	err = p.QueryRow(query).Scan(
		&cartSession,
	)
	if err != nil {
		return "", err
	}
	return cartSession, nil
}

// ChangeQuantity ..
func (p *PostgresDatastore) ChangeQuantity(variantID string, cartSession string, newQuantety string) (err error) {
	query := fmt.Sprintf(`SELECT * FROM cart.change_product(%s, %s, %s);`, variantID, cartSession, newQuantety)
	err = p.QueryRow(query).Scan(nil)
	if err != nil {
		return err
	}
	return nil
}

// RemoveProduct ..
func (p *PostgresDatastore) RemoveProduct(cartSession string, variant *models.Variant) (err error) {
	query := fmt.Sprintf(`SELECT * FROM cart.remove_product(%v);`, variant)
	err = p.QueryRow(query).Scan(nil)
	if err != nil {
		return err
	}
	return nil
}

// GetCart ..
func (p *PostgresDatastore) GetCart(cartSession string) (cartItems []*models.Variant, err error) {
	query := fmt.Sprintf(`SELECT * FROM cart.get_cart(%s);`, cartSession)

	variant := &models.Variant{}

	rows, err := p.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(
			variant.VariantID,
			variant.Attributes,
			variant.Description,
			variant.Images,
			variant.Name,
			variant.PriceOverride,
			variant.ProductID,
			variant.Quantity,
		)
		if err != nil {
			return nil, err
		}

		cartItems = append(cartItems, variant)
	}

	if err != nil {
		return cartItems, err
	}

	return cartItems, nil
}
