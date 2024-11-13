package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TeamMeta struct {
	bun.BaseModel `bun:"table:teams_meta"`

	ID *uuid.UUID `bun:"id,pk,type:uuid"`

	TeamID *uuid.UUID `bun:"team_id"`
	UserID string     `bun:"user_id"`
}
