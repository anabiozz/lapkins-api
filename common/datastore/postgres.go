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
func (p *PostgresDatastore) GetProducts(subjectURL string) (products []models.Product, err error) {
	fmt.Println(subjectURL)
	query := fmt.Sprintf(`SELECT * FROM new_products.get_products('%s');`, subjectURL)
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

// GetVariation ..
func (p *PostgresDatastore) GetVariation(variationID, sizeOptionID string) (*models.Variation, error) {
	query := fmt.Sprintf(`SELECT * FROM new_products.get_variation(%s, %s);`, variationID, sizeOptionID)

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
		&variation.Size,
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
func (p *PostgresDatastore) AddProduct(variationID int, сartSession string) (err error) {
	query := fmt.Sprintf(`SELECT * FROM cart.add_product(%d, '%s');`, variationID, сartSession)
	_, err = p.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// ChangeQuantity ..
func (p *PostgresDatastore) ChangeQuantity(variationID string, cartSession string, newQuantety string) (err error) {
	query := fmt.Sprintf(`SELECT * FROM cart.change_product(%s, %s, %s);`, variationID, cartSession, newQuantety)
	err = p.QueryRow(query).Scan(nil)
	if err != nil {
		return err
	}
	return nil
}

// RemoveProduct ..
func (p *PostgresDatastore) RemoveProduct(cartSession string, variation *models.Variation) (err error) {
	query := fmt.Sprintf(`SELECT * FROM cart.remove_product(%v);`, variation)
	err = p.QueryRow(query).Scan(nil)
	if err != nil {
		return err
	}
	return nil
}

// GetCart ..
func (p *PostgresDatastore) GetCart(cartSession string) (cartItems []*models.Variation, err error) {
	query := fmt.Sprintf(`SELECT * FROM cart.get_cart(%s);`, cartSession)

	variation := &models.Variation{}

	rows, err := p.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(
			variation.ID,
			variation.Attributes,
			variation.Description,
			variation.Images,
			variation.Name,
			variation.Price,
			variation.ProductID,
		)
		if err != nil {
			return nil, err
		}

		cartItems = append(cartItems, variation)
	}

	if err != nil {
		return cartItems, err
	}

	return cartItems, nil
}
