package repository

import (
	"context"
	"errors"

	"github.com/alvinmatias69/wedding_invitation/internal/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) BeginTrx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *Repository) GetByJwtToken(ctx context.Context, token string) (entities.Token, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, jwt_token, steam_token, claimed_at FROM tokens WHERE jwt_token = $1;", token)
	if err != nil {
		return entities.Token{}, err
	}

	payload, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Token])
	if errors.Is(err, pgx.ErrNoRows) {
		return entities.Token{}, errors.New("not found")
	}

	return payload, err
}

func (r *Repository) FindOneUnclaimed(ctx context.Context, trx pgx.Tx) (entities.Token, error) {
	rows, err := trx.Query(ctx, "SELECT id, jwt_token, steam_token, claimed_at FROM tokens WHERE jwt_token = NULL;")
	if err != nil {
		return entities.Token{}, err
	}

	payload, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Token])
	if errors.Is(err, pgx.ErrNoRows) {
		return entities.Token{}, errors.New("not found")
	}

	return payload, err
}

func (r *Repository) Claim(ctx context.Context, trx pgx.Tx, id uint64, jwtToken string) error {
	_, err := trx.Exec(ctx, "UPDATE tokens SET claimed_at = NOW(), jwt_token = $1 WHERE id = $2;", jwtToken, id)
	return err
}
