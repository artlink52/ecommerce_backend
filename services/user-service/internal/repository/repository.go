package repository

import (
	"context"
	"errors"
	"fmt"

	domainerrors "github.com/artlink52/ecommerce_backend/services/user-service/internal/domain/errors"
	"github.com/artlink52/ecommerce_backend/services/user-service/internal/domain/models"
	pgxpool "github.com/artlink52/ecommerce_backend/services/user-service/internal/repository/pgx"
	pgxlib "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const pgUniqueViolation = "23505"

type UserRepository struct {
	pool pgxpool.Pool
}

func NewUserRepository(pool pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) SaveUser(ctx context.Context, email string, passwordHash []byte) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	const query = `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`

	var id int64
	err := r.pool.QueryRow(ctx, query, email, passwordHash).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return 0, domainerrors.ErrUserExists
		}
		return 0, fmt.Errorf("save user: %w", err)
	}
	return id, nil
}

func (r *UserRepository) User(ctx context.Context, email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	const query = `SELECT id, email, password_hash FROM users WHERE email = $1`

	var user models.User
	err := r.pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, pgxlib.ErrNoRows) {
			return models.User{}, domainerrors.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}
