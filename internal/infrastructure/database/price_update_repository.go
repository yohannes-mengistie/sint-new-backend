package database

import (
	"context"
	"fmt"

	"github.com/yohannes/sint-new-backend/internal/domain"
)

func (database *Postgres) CreatePriceUpdate(
	ctx context.Context,
	user domain.User,
	update domain.PriceUpdate,
) error {
	tx, err := database.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin price update transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var userID int64
	err = tx.QueryRow(
		ctx,
		`SELECT id FROM users WHERE firebase_uid = $1`,
		user.FirebaseUID,
	).Scan(&userID)
	if err != nil {
		return fmt.Errorf("find submitting user: %w", err)
	}

	var productID int64
	err = tx.QueryRow(
		ctx,
		`INSERT INTO products (name)
		 VALUES ($1)
		 ON CONFLICT DO NOTHING
		 RETURNING id`,
		update.ProductName,
	).Scan(&productID)
	if err != nil {
		err = tx.QueryRow(
			ctx,
			`SELECT id FROM products WHERE name = $1 ORDER BY id LIMIT 1`,
			update.ProductName,
		).Scan(&productID)
		if err != nil {
			return fmt.Errorf("find product: %w", err)
		}
	}

	var storeID int64
	err = tx.QueryRow(
		ctx,
		`INSERT INTO stores (name)
		 VALUES ($1)
		 ON CONFLICT DO NOTHING
		 RETURNING id`,
		update.StoreName,
	).Scan(&storeID)
	if err != nil {
		err = tx.QueryRow(
			ctx,
			`SELECT id FROM stores WHERE name = $1 ORDER BY id LIMIT 1`,
			update.StoreName,
		).Scan(&storeID)
		if err != nil {
			return fmt.Errorf("find store: %w", err)
		}
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO price_updates
		 (product_id, store_id, submitted_by, price, photo_url, notes, status, observed_at)
		 VALUES ($1, $2, $3, $4, NULLIF($5, ''), NULLIF($6, ''), 'pending', NOW())`,
		productID,
		storeID,
		userID,
		update.Price,
		update.PhotoURL,
		update.Notes,
	)
	if err != nil {
		return fmt.Errorf("create price update: %w", err)
	}
	return tx.Commit(ctx)
}
