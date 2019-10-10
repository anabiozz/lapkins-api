package datastore

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/anabiozz/lapkin-project/lapkin-api/models"
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

	logger.Info(db.DBinfo)

	connection, err := sql.Open("postgres", db.DBinfo)
	if err != nil {
		return nil, err
	}
	return &PostgresDatastore{
		DB: connection,
	}, nil
}

// GetProducts ..
func (p *PostgresDatastore) GetProducts(productsID string) (products []models.Product, err error) {
	id, err := strconv.Atoi(productsID)
	query := fmt.Sprintf(`SELECT * FROM products.get_products(%d);`, id)
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
			&product.Price,
			&product.Size)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

// GetVariant ..
func (p *PostgresDatastore) GetVariant(variantID, size string) (*models.Variant, error) {
	query := fmt.Sprintf(`SELECT * FROM products.get_variant(%s, '%s');`, variantID, size)

	variant := &models.Variant{}

	err := p.QueryRow(query).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.Name,
		&variant.Description,
		&variant.PriceOverride,
		&variant.Attributes,
		pq.Array(&variant.Sizes),
		&variant.Size,
		pq.Array(&variant.Images),
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
