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
			// &product.ProductVariantName,
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			pq.Array(&product.Images))
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

// GetProductByID ..
func (p *PostgresDatastore) GetProductByID(productID string) (product *models.ProductVariant, err error) {
	id, err := strconv.Atoi(productID)
	query := fmt.Sprintf(`SELECT * FROM products.get_product_by_id(%d);`, id)

	product = &models.ProductVariant{}

	err = p.QueryRow(query).Scan(
		&product.ID,
		&product.ProductID,
		&product.Name,
		&product.Description,
		&product.PriceOverride,
		&product.Attributes,
		pq.Array(&product.Sizes),
		pq.Array(&product.Images))
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductVariantByID ..
func (p *PostgresDatastore) GetProductVariantByID(productVariantID, size string) (product *models.ProductVariant, err error) {
	id, err := strconv.Atoi(productVariantID)
	query := fmt.Sprintf(`SELECT * FROM products.get_product_variant_by_id(%d, '%s');`, id, size)

	product = &models.ProductVariant{}

	err = p.QueryRow(query).Scan(
		&product.ID,
		&product.ProductID,
		&product.Name,
		&product.Description,
		&product.PriceOverride,
		&product.Attributes,
		pq.Array(&product.Sizes),
		pq.Array(&product.Images))
	if err != nil {
		return nil, err
	}

	return product, nil
}

// CloseDB ..
func (p *PostgresDatastore) CloseDB() {
	p.DB.Close()
}
