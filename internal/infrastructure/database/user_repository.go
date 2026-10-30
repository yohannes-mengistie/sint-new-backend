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
