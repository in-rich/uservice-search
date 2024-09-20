package dao

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/uptrace/bun"
)

type SearchNotesRepository interface {
	SearchNotes(ctx context.Context, limit int, offset int, rawQuery string, authorID string) ([]*entities.Note, error)
}

type searchNotesRepositoryImpl struct {
	db bun.IDB
}

func (r *searchNotesRepositoryImpl) SearchNotes(ctx context.Context, limit int, offset int, rawQuery string, authorID string) ([]*entities.Note, error) {
	notes := make([]*entities.Note, 0)

	var query *bun.SelectQuery

	if rawQuery != "" {
		formattedSearchQuery := r.db.NewSelect().
			ColumnExpr(`to_tsquery('english', string_agg(lexeme || ':*', ' & ' ORDER BY POSITIONS)) AS query`).
			TableExpr(`unnest(to_tsvector('english', unaccent(?)))`, rawQuery)

		query = r.db.NewSelect().Model(&notes).
			TableExpr("(?) AS search", formattedSearchQuery).
			Where("(note.content || note.target_name) @@ search.query").
			Where("author_id = ?", authorID).
			OrderExpr("ts_rank_cd((note.content || note.target_name), search.query) DESC").
			OrderExpr("updated_at DESC").
			Limit(limit).Offset(offset)
	} else {
		query = r.db.NewSelect().Model(&notes).
			Where("author_id = ?", authorID).
			Order("updated_at DESC").
			Limit(limit).Offset(offset)
	}

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}

	if len(notes) == 0 {
		return nil, ErrNoteNotFound
	}

	return notes, nil
}

func NewSearchNotesRepository(db bun.IDB) SearchNotesRepository {
	return &searchNotesRepositoryImpl{
		db: db,
	}
}
