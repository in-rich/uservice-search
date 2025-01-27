package dao

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/uptrace/bun"
)

type SearchRemindersRepository interface {
	SearchReminders(ctx context.Context, authorID string, rawQuery string, limit int, offset int) ([]*entities.Reminder, error)
}

type searchRemindersRepositoryImpl struct {
	db bun.IDB
}

func (r *searchRemindersRepositoryImpl) SearchReminders(ctx context.Context, authorID string, rawQuery string, limit int, offset int) ([]*entities.Reminder, error) {
	reminders := make([]*entities.Reminder, 0)

	query := r.db.NewSelect().Model(&reminders).
		Where("author_id = ?", authorID).
		Limit(limit).Offset(offset)

	if rawQuery != "" {
		formattedSearchQuery := r.db.NewSelect().
			ColumnExpr(`to_tsquery('english', string_agg(lexeme || ':*', ' & ' ORDER BY POSITIONS)) AS query`).
			TableExpr(`unnest(to_tsvector('english', unaccent(?)))`, rawQuery)

		query = query.
			TableExpr("(?) AS search", formattedSearchQuery).
			Where("(reminder.content || reminder.target_name) @@ search.query").
			OrderExpr("ts_rank_cd((reminder.content || reminder.target_name), search.query) DESC")
	}

	// Add default orderBy last, so it does not take precedence over the score.
	query = query.Order("expired_at DESC")

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}

	return reminders, nil
}

func NewSearchRemindersRepository(db bun.IDB) SearchRemindersRepository {
	return &searchRemindersRepositoryImpl{
		db: db,
	}
}
