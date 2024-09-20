package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

const searchMessageReturning = "team_id, message_id, id"

type Message struct {
	bun.BaseModel `bun:"table:messages"`

	ID *uuid.UUID `bun:"id,pk,type:uuid"`

	TeamID    string `bun:"team_id"`
	MessageID string `bun:"message_id"`

	Content    string `bun:"content"`
	TargetName string `bun:"target_name"`

	UpdatedAt *time.Time `bun:"updated_at,notnull"`
}

func (m *Message) BeforeCreate(query *bun.InsertQuery) *bun.InsertQuery {
	return query.
		Value("content", Vectorize("A"), m.Content).
		Value("target_name", Vectorize("A"), m.TargetName).
		Returning(searchMessageReturning)
}

func (m *Message) BeforeUpdate(query *bun.UpdateQuery) *bun.UpdateQuery {
	return query.
		Column("content", "target_name").
		Value("content", Vectorize("A"), m.Content).
		Value("target_name", Vectorize("A"), m.TargetName).
		Returning(searchMessageReturning)
}
