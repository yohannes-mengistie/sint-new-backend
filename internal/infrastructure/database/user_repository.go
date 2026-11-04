package database

import (
	"context"

	"github.com/yohannes/sint-new-backend/internal/domain"
)

func (database *Postgres) Upsert(ctx context.Context, user domain.User) (domain.User, error) {
	_, err := database.pool.Exec(
		ctx,
		`INSERT INTO users (firebase_uid, email, display_name)
		 VALUES ($1, NULLIF($2, ''), NULLIF($3, ''))
		 ON CONFLICT (firebase_uid)
		 DO UPDATE SET email = EXCLUDED.email,
		               display_name = EXCLUDED.display_name,
		               updated_at = NOW()`,
		user.FirebaseUID,
		user.Email,
		user.DisplayName,
	)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (database *Postgres) ListProductIDs(
	ctx context.Context,
	user domain.User,
) ([]int64, error) {
	rows, err := database.pool.Query(ctx, `
		SELECT f.product_id
		FROM favorites f
		JOIN users u ON u.id = f.user_id
		WHERE u.firebase_uid = $1
		ORDER BY f.created_at DESC`, user.FirebaseUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (database *Postgres) Add(
	ctx context.Context,
	user domain.User,
	productID int64,
) error {
	_, err := database.pool.Exec(ctx, `
		INSERT INTO favorites (user_id, product_id)
		SELECT u.id, $2
		FROM users u
		JOIN products p ON p.id = $2
		WHERE u.firebase_uid = $1
		ON CONFLICT (user_id, product_id) DO NOTHING`,
		user.FirebaseUID,
		productID,
	)
	return err
}

func (database *Postgres) Remove(
	ctx context.Context,
	user domain.User,
	productID int64,
) error {
	_, err := database.pool.Exec(ctx, `
		DELETE FROM favorites
		WHERE user_id = (
			SELECT id FROM users WHERE firebase_uid = $1
		) AND product_id = $2`,
		user.FirebaseUID,
		productID,
	)
	return err
}

func (database *Postgres) ListAlerts(
	ctx context.Context,
	user domain.User,
) ([]domain.Alert, error) {
	rows, err := database.pool.Query(ctx, `
		SELECT a.id, a.type, a.title, a.body, a.created_at, a.read_at
		FROM alerts a
		JOIN users u ON u.id = a.user_id
		WHERE u.firebase_uid = $1
		ORDER BY a.created_at DESC`, user.FirebaseUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	alerts := make([]domain.Alert, 0)
	for rows.Next() {
		var alert domain.Alert
		if err := rows.Scan(
			&alert.ID,
			&alert.Type,
			&alert.Title,
			&alert.Body,
			&alert.CreatedAt,
			&alert.ReadAt,
		); err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}
	return alerts, rows.Err()
}

func (database *Postgres) DeleteAlert(
	ctx context.Context,
	user domain.User,
	alertID int64,
) error {
	_, err := database.pool.Exec(ctx, `
		DELETE FROM alerts
		WHERE id = $1
		  AND user_id = (
			SELECT id FROM users WHERE firebase_uid = $2
		  )`, alertID, user.FirebaseUID)
	return err
}
