package database

import (
	"context"

	"github.com/yohannes/sint-new-backend/internal/domain"
)

func (database *Postgres) List(ctx context.Context) ([]domain.Product, error) {
	rows, err := database.pool.Query(
		ctx,
		`SELECT p.id, p.name, COALESCE(p.variant, ''), COALESCE(p.size, ''),
		        COALESCE(p.barcode, ''), COALESCE(p.image_url, ''),
		        COALESCE(latest.price, 0), COALESCE(latest.store_name, '')
		 FROM products p
		 LEFT JOIN LATERAL (
		   SELECT pu.price, s.name AS store_name
		   FROM price_updates pu
		   JOIN stores s ON s.id = pu.store_id
		   WHERE pu.product_id = p.id AND pu.status = 'approved'
		   ORDER BY pu.observed_at DESC, pu.created_at DESC
		   LIMIT 1
		 ) latest ON TRUE
		 ORDER BY name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]domain.Product, 0)
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Variant,
			&product.Size,
			&product.Barcode,
			&product.ImageURL,
			&product.Price,
			&product.Store,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}
