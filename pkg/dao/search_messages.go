package dao

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/uptrace/bun"
)

type SearchMessagesRepository interface {
	SearchMessages(ctx context.Context, teamID string, rawQuery string, limit int, offset int) ([]*entities.Message, error)
}

type searchMessagesRepositoryImpl struct {
	db bun.IDB
}

func (r *searchMessagesRepositoryImpl) SearchMessages(ctx context.Context, teamID string, rawQuery string, limit int, offset int) ([]*entities.Message, error) {
	messages := make([]*entities.Message, 0)

	query := r.db.NewSelect().Model(&messages).
		Where("team_id = ?", teamID).
		Limit(limit).Offset(offset)

	if rawQuery != "" {
		formattedSearchQuery := r.db.NewSelect().
			ColumnExpr(`to_tsquery('english', string_agg(lexeme || ':*', ' & ' ORDER BY POSITIONS)) AS query`).
			TableExpr(`unnest(to_tsvector('english', unaccent(?)))`, rawQuery)

		query = query.
			TableExpr("(?) AS search", formattedSearchQuery).
			Where("(message.content || message.target_name) @@ search.query").
			OrderExpr("ts_rank_cd((message.content || message.target_name), search.query) DESC")
	}

	// Add default orderBy last, so it does not take precedence over the score.
	query = query.Order("updated_at DESC")

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func NewSearchMessagesRepository(db bun.IDB) SearchMessagesRepository {
	return &searchMessagesRepositoryImpl{
		db: db,
	}
}
