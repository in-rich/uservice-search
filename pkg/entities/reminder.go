package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

const searchRemindersReturning = "author_id, reminder_id, id"

type Reminder struct {
	bun.BaseModel `bun:"table:reminders"`

	ID *uuid.UUID `bun:"id,pk,type:uuid"`

	AuthorID   string `bun:"author_id"`
	ReminderID string `bun:"reminder_id"`

	Content    string `bun:"content"`
	TargetName string `bun:"target_name"`

	UpdatedAt *time.Time `bun:"updated_at,notnull"`
	ExpiredAt *time.Time `bun:"expired_at,notnull"`
}

func (n *Reminder) BeforeCreate(query *bun.InsertQuery) *bun.InsertQuery {
	return query.
		Value("content", Vectorize("A"), n.Content).
		Value("target_name", Vectorize("A"), n.TargetName).
		Returning(searchRemindersReturning)
}

func (n *Reminder) BeforeUpdate(query *bun.UpdateQuery) *bun.UpdateQuery {
	return query.
		Column("content", "target_name", "updated_at", "expired_at").
		Value("content", Vectorize("A"), n.Content).
		Value("target_name", Vectorize("A"), n.TargetName).
		Returning(searchRemindersReturning)
}
