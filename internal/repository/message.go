package repository

import (
	"context"
	"errors"

	"github.com/alvinmatias69/wedding_invitation/internal/constant"
	"github.com/alvinmatias69/wedding_invitation/internal/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	pool *pgxpool.Pool
}

func NewMessageRepository(pool *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{
		pool: pool,
	}
}

func (m *MessageRepository) Insert(ctx context.Context, msg entities.Message) error {
	_, err := m.pool.Exec(ctx, "INSERT INTO messages(sender_name, content, created_at) VALUES ($1, $2, NOW());", msg.SenderName, msg.Content)
	return err
}

func (m *MessageRepository) Get(ctx context.Context, limit uint64, offset uint64) ([]entities.Message, error) {
	rows, err := m.pool.Query(ctx, "SELECT sender_name, content, created_at FROM messages ORDER BY created_at DESC LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return nil, err
	}

	payload, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.Message])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, constant.ErrNotFound
	}

	return payload, err
}
