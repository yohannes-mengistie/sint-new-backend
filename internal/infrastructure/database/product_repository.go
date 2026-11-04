package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/yohannes/sint-new-backend/internal/application"
	"github.com/yohannes/sint-new-backend/internal/domain"
)

func (database *Postgres) List(ctx context.Context) ([]domain.Product, error) {
	rows, err := database.pool.Query(ctx, `
		SELECT p.id, p.name, COALESCE(p.variant, ''), COALESCE(p.size, ''),
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
		ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]domain.Product, 0)
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(
			&product.ID, &product.Name, &product.Variant, &product.Size,
			&product.Barcode, &product.ImageURL, &product.Price, &product.Store,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}

func (database *Postgres) ListStores(ctx context.Context) ([]domain.Store, error) {
	rows, err := database.pool.Query(ctx, `
		SELECT id, name, COALESCE(address, ''), COALESCE(latitude, 0),
		       COALESCE(longitude, 0)
		FROM stores
		ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stores := make([]domain.Store, 0)
	for rows.Next() {
		var store domain.Store
		if err := rows.Scan(
			&store.ID,
			&store.Name,
			&store.Address,
			&store.Latitude,
			&store.Longitude,
		); err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}
	return stores, rows.Err()
}

func (database *Postgres) GetDetail(ctx context.Context, id int64) (domain.ProductDetail, error) {
	var detail domain.ProductDetail
	err := database.pool.QueryRow(ctx, `
		SELECT id, name, COALESCE(variant, ''), COALESCE(size, ''),
		       COALESCE(barcode, ''), COALESCE(image_url, '')
		FROM products WHERE id = $1`, id).Scan(
		&detail.ID, &detail.Name, &detail.Variant, &detail.Size,
		&detail.Barcode, &detail.ImageURL,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ProductDetail{}, application.ErrProductNotFound
	}
	if err != nil {
		return domain.ProductDetail{}, err
	}

	detail.Description = detail.Name + " " + detail.Size + ", available at stores near you."
	detail.About = detail.Description + " Prices are shared by the community and updated in real time."
	detail.Offers = make([]domain.StoreOffer, 0)
	rows, err := database.pool.Query(ctx, `
		SELECT s.id, s.name, COALESCE(s.address, ''), COALESCE(s.latitude, 0),
		       COALESCE(s.longitude, 0), current_price.price,
		       COALESCE(previous_price.price, current_price.price)
		FROM stores s
		JOIN LATERAL (
		  SELECT pu.price
		  FROM price_updates pu
		  WHERE pu.product_id = $1 AND pu.store_id = s.id AND pu.status = 'approved'
		  ORDER BY pu.observed_at DESC, pu.created_at DESC LIMIT 1
		) current_price ON TRUE
		LEFT JOIN LATERAL (
		  SELECT pu.price
		  FROM price_updates pu
		  WHERE pu.product_id = $1 AND pu.store_id = s.id AND pu.status = 'approved'
		  ORDER BY pu.observed_at DESC, pu.created_at DESC OFFSET 1 LIMIT 1
		) previous_price ON TRUE
		ORDER BY current_price.price ASC`, id)
	if err != nil {
		return domain.ProductDetail{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var offer domain.StoreOffer
		if err := rows.Scan(
			&offer.Store.ID, &offer.Store.Name, &offer.Store.Address,
			&offer.Store.Latitude, &offer.Store.Longitude,
			&offer.Price, &offer.OldPrice,
		); err != nil {
			return domain.ProductDetail{}, err
		}
		offer.InStock = true
		detail.Offers = append(detail.Offers, offer)
	}
	if err := rows.Err(); err != nil {
		return domain.ProductDetail{}, err
	}

	detail.History = make([]domain.PricePoint, 0)
	historyRows, err := database.pool.Query(ctx, `
		SELECT TO_CHAR(day, 'Dy'), ROUND(AVG(price), 2)
		FROM (
		  SELECT date_trunc('day', observed_at)::date AS day, price
		  FROM price_updates
		  WHERE product_id = $1 AND status = 'approved'
		    AND observed_at >= CURRENT_DATE - INTERVAL '6 days'
		) prices
		GROUP BY day ORDER BY day`, id)
	if err != nil {
		return domain.ProductDetail{}, err
	}
	defer historyRows.Close()
	for historyRows.Next() {
		var point domain.PricePoint
		if err := historyRows.Scan(&point.Day, &point.Price); err != nil {
			return domain.ProductDetail{}, err
		}
		detail.History = append(detail.History, point)
	}
	return detail, historyRows.Err()
}
