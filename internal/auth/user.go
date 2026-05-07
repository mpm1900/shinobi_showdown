package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"shinobi_showdown/internal/db"
)

func GetUser(ctx context.Context, queries *db.Queries, email string) (*db.User, bool, error) {
	user, err := queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &user, true, nil
}

func CreateUser(ctx context.Context, queries *db.Queries, username, email, password string) (*db.User, error) {
	hashed, salt, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Username: username,
		Email:    email,
		Password: hashed,
		Salt:     salt,
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}
