package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hareshkhan01/PollyRoute/internals/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, name string, email string, password_hash string) (string, error)
	UpdateRefreshToken(ctx context.Context, userId string, token *string, expiresAt *time.Time) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByUserId(ctx context.Context, userId string) (*domain.User, error)
	FindByRefreshToken(ctx context.Context, token string) (*domain.User, error)
	UpdateUser(ctx context.Context, id string, name string, email string, password_hash string) error
	DeleteUserById(ctx context.Context, userId string) error
	DeleteUserByEmail(ctx context.Context, email string) error
}

type userRepository struct {
	dbPool *pgxpool.Pool
}

func NewUserRepository(dbPool *pgxpool.Pool) UserRepository {
	return &userRepository{
		dbPool: dbPool,
	}
}

func (ur *userRepository) CreateUser(
	ctx context.Context,
	name string,
	email string,
	password_hash string,
) (string, error) {
	query := `INSERT INTO users(name,email,password_hash) VALUES ($1,$2,$3) RETURNING id`

	var insertId string
	if err := ur.dbPool.QueryRow(ctx, query, name, email, password_hash).Scan(&insertId); err != nil {
		return "", fmt.Errorf("Failed to Create User: %w", err)
	}
	return insertId, nil
}

func (ur *userRepository) UpdateRefreshToken(
	ctx context.Context,
	userId string,
	token *string,
	expiresAt *time.Time,
) error {
	query := `UPDATE users SET refresh_token=$1,refresh_token_expires_at=$2 WHERE id=$3`

	commandTag, err := ur.dbPool.Exec(ctx, query, token, expiresAt, userId)
	if err != nil {
		return fmt.Errorf("Failed to update token: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (ur *userRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {
	query := `SELECT id,name,email,password_hash,refresh_token,refresh_token_expires_at,created_at FROM users WHERE email=$1`

	var user domain.User

	err := ur.dbPool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.RefreshToken,
		&user.RefreshTokenExpiresAt,
		&user.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return &user, nil

}

func (ur *userRepository) FindByUserId(
	ctx context.Context,
	userId string,
) (*domain.User, error) {
	query := `SELECT id,name,email,password_hash,refresh_token,refresh_token_expires_at,created_at FROM users WHERE id=$1`

	var user domain.User

	err := ur.dbPool.QueryRow(ctx, query, userId).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.RefreshToken,
		&user.RefreshTokenExpiresAt,
		&user.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return &user, nil

}

func (ur *userRepository) UpdateUser(
	ctx context.Context,
	id string,
	name string,
	email string,
	password_hash string,
) error {
	query := `UPDATE users SET name=$1,email=$2,password_hash=$3 WHERE id=$4`
	commandTag, err := ur.dbPool.Exec(ctx, query, name, email, password_hash, id)
	if err != nil {
		return fmt.Errorf("Repository: Failed to update user -  %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (ur *userRepository) FindByRefreshToken(ctx context.Context, token string) (*domain.User, error) {
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("Token is Empty!")
	}
	query := `SELECT id,name,email,password_hash,refresh_token,refresh_token_expires_at,created_at FROM users WHERE refresh_token=$1`

	var user domain.User
	err := ur.dbPool.QueryRow(ctx, query, token).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.RefreshToken,
		&user.RefreshTokenExpiresAt,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) DeleteUserById(ctx context.Context, userId string) error {
	query := `DELETE FROM users WHERE id=$1`
	commandTag, err := ur.dbPool.Exec(ctx, query, userId)

	if err != nil {
		return fmt.Errorf("Failed to delete %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (ur *userRepository) DeleteUserByEmail(ctx context.Context, email string) error {
	query := `DELETE FROM users WHERE email=$1`
	commandTag, err := ur.dbPool.Exec(ctx, query, email)

	if err != nil {
		return fmt.Errorf("Failed to delete %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
