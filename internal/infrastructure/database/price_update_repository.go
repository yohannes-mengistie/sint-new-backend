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

	_, err = tx.Exec(
		ctx,
		`INSERT INTO price_updates
		 (product_id, store_id, submitted_by, price, photo_url, notes, status, observed_at)
		 VALUES ($1, $2, $3, $4, NULLIF($5, ''), NULLIF($6, ''), 'pending', NOW())`,
		update.ProductID,
		update.StoreID,
		userID,
		update.Price,
		update.PhotoURL,
		update.Notes,
	)
	if err != nil {
		return fmt.Errorf("create price update: %w", err)
	}
	_, err = tx.Exec(
		ctx,
		`INSERT INTO alerts (user_id, type, title, body)
		 VALUES ($1, 'submission', 'Price update submitted',
		         'Your price update is waiting for review.')`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("create submission alert: %w", err)
	}
	return tx.Commit(ctx)
}

func (database *Postgres) ListPriceSubmissions(
	ctx context.Context,
	user domain.User,
) ([]domain.PriceSubmission, error) {
	rows, err := database.pool.Query(ctx, `
		SELECT pu.id, p.name, s.name, pu.price, pu.status,
		       COALESCE(pu.photo_url, ''), COALESCE(pu.notes, ''),
		       pu.observed_at, pu.created_at
		FROM price_updates pu
		JOIN products p ON p.id = pu.product_id
		JOIN stores s ON s.id = pu.store_id
		JOIN users u ON u.id = pu.submitted_by
		WHERE u.firebase_uid = $1
		ORDER BY pu.created_at DESC`, user.FirebaseUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	submissions := make([]domain.PriceSubmission, 0)
	for rows.Next() {
		var submission domain.PriceSubmission
		if err := rows.Scan(
			&submission.ID,
			&submission.ProductName,
			&submission.StoreName,
			&submission.Price,
			&submission.Status,
			&submission.PhotoURL,
			&submission.Notes,
			&submission.ObservedAt,
			&submission.CreatedAt,
		); err != nil {
			return nil, err
		}
		submissions = append(submissions, submission)
	}
	return submissions, rows.Err()
}
