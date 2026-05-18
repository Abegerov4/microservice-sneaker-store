package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"sneaker-store/ai-service/internal/model"
)

type ChatRepository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) SaveMessage(ctx context.Context, msg *model.ChatMessage) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO chat_history (id, session_id, user_id, role, content, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		msg.ID, msg.SessionID, msg.UserID, msg.Role, msg.Content, msg.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("save message: %w", err)
	}
	return nil
}

func (r *ChatRepository) GetHistory(ctx context.Context, sessionID string, limit int) ([]*model.ChatMessage, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, session_id, user_id, role, content, created_at
		 FROM chat_history
		 WHERE session_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2`,
		sessionID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}
	defer rows.Close()

	var msgs []*model.ChatMessage
	for rows.Next() {
		m := &model.ChatMessage{}
		if err := rows.Scan(&m.ID, &m.SessionID, &m.UserID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		msgs = append(msgs, m)
	}

	// reverse to chronological order
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs, nil
}
