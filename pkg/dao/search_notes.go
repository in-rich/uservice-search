package dao

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/uptrace/bun"
)

type SearchNotesRepository interface {
	SearchNotes(ctx context.Context, authorID string, rawQuery string, limit int, offset int) ([]*entities.Note, error)
}

type searchNotesRepositoryImpl struct {
	db bun.IDB
}

func (r *searchNotesRepositoryImpl) SearchNotes(ctx context.Context, authorID string, rawQuery string, limit int, offset int) ([]*entities.Note, error) {
	notes := make([]*entities.Note, 0)

	query := r.db.NewSelect().Model(&notes).
		Where("author_id = ?", authorID).
		Limit(limit).Offset(offset)

	if rawQuery != "" {
		formattedSearchQuery := r.db.NewSelect().
			ColumnExpr(`to_tsquery('english', string_agg(lexeme || ':*', ' & ' ORDER BY POSITIONS)) AS query`).
			TableExpr(`unnest(to_tsvector('english', unaccent(?)))`, rawQuery)

		query = query.
			TableExpr("(?) AS search", formattedSearchQuery).
			Where("(note.content || note.target_name) @@ search.query").
			OrderExpr("ts_rank_cd((note.content || note.target_name), search.query) DESC")
	}

	// Add default orderBy last, so it does not take precedence over the score.
	query = query.Order("updated_at DESC")

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func NewSearchNotesRepository(db bun.IDB) SearchNotesRepository {
	return &searchNotesRepositoryImpl{
		db: db,
	}
}
