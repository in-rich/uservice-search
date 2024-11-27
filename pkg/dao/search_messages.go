package dao

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

type SearchMessageData struct {
	UserID           string
	TeamID           string
	RawQuery         string
	Limit            int
	Offset           int
	OneMessageByTeam bool
}

type SearchMessagesRepository interface {
	SearchMessages(ctx context.Context, data *SearchMessageData) ([]*entities.Message, error)
}

type searchMessagesRepositoryImpl struct {
	db bun.IDB
}

type searchMessageFilter struct {
	query  string
	params []any
}

type searchMessageTableExpr struct {
	query  string
	params []any
}

func (r *searchMessagesRepositoryImpl) SearchMessages(ctx context.Context, data *SearchMessageData) ([]*entities.Message, error) {
	// Build the filter for the final set of messages.
	var filters []searchMessageFilter
	var orderBy []string
	var tableExprs []searchMessageTableExpr

	// If we have a UserID we get all of user's teams and filter the messages with the teams.
	if data.UserID != "" {
		teams := make([]entities.TeamMeta, 0)
		teamsIDs := make([]*uuid.UUID, 0)
		teamQuery := r.db.NewSelect().Model(&teams).Column().Where("user_id = ?", data.UserID)

		if err := teamQuery.Scan(ctx); err != nil {
			fmt.Println("coucou")
			return nil, fmt.Errorf("get team ids: %w", err)
		}

		for _, team := range teams {
			teamsIDs = append(teamsIDs, team.TeamID)
		}

		filters = append(filters, searchMessageFilter{
			query:  "team_id IN (?)",
			params: []any{bun.In(teamsIDs)},
		})
	}

	// If we have a TeamID we search all of the messages on this team.
	if data.TeamID != "" {
		filters = append(filters, searchMessageFilter{
			query:  "team_id = ?",
			params: []any{data.TeamID},
		})
	}

	// If we have a RawQuery we filter the message with the query.
	if data.RawQuery != "" {
		formattedSearchQuery := r.db.NewSelect().
			ColumnExpr(`to_tsquery('english', string_agg(lexeme || ':*', ' & ' ORDER BY POSITIONS)) AS messageQuery`).
			TableExpr(`unnest(to_tsvector('english', unaccent(?)))`, data.RawQuery)

		filters = append(filters, searchMessageFilter{
			query: "(content || target_name) @@ messageQuery",
		})

		tableExprs = append(tableExprs, searchMessageTableExpr{
			query:  "(?) AS search",
			params: []any{formattedSearchQuery},
		})

		orderBy = append(orderBy, "ts_rank_cd((content || target_name), messageQuery) DESC")
	}

	orderBy = append(orderBy, "updated_at DESC")

	messages := make([]*entities.Message, 0)

	query := r.db.NewSelect().Model(&messages)

	// We separate the creation of the query whether we want one or more message.
	query = lo.
		IfF(data.OneMessageByTeam, func() *bun.SelectQuery {
			subQuery := r.db.NewSelect().
				TableExpr("messages as unique_messages").
				ColumnExpr("unique_messages.*").
				GroupExpr("team_id, id, target_name").
				Order("team_id DESC", "target_name DESC").
				DistinctOn("team_id, target_name").
				Order("updated_at DESC")

			if data.RawQuery != "" {
				subQuery = subQuery.ColumnExpr("search.messageQuery as messageQuery")
				subQuery = subQuery.GroupExpr("messageQuery")
			}

			for _, filter := range filters {
				subQuery = subQuery.Where(filter.query, filter.params...)
			}

			for _, tableExpr := range tableExprs {
				subQuery = subQuery.TableExpr(tableExpr.query, tableExpr.params...)
			}

			res := query.ModelTableExpr("(?) as message", subQuery)

			for _, order := range orderBy {
				res = res.OrderExpr(order)
			}

			return res
		}).
		ElseF(func() *bun.SelectQuery {
			res := query

			for _, filter := range filters {
				res = res.Where(filter.query, filter.params...)
			}

			for _, tableExpr := range tableExprs {
				res = res.TableExpr(tableExpr.query, tableExpr.params...)
			}

			for _, order := range orderBy {
				res = res.OrderExpr(order)
			}

			return res
		})

	if err := query.Scan(ctx); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	return messages, nil
}

func NewSearchMessagesRepository(db bun.IDB) SearchMessagesRepository {
	return &searchMessagesRepositoryImpl{
		db: db,
	}
}
