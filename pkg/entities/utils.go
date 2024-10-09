package entities

import (
	"fmt"
	"github.com/uptrace/bun"
)

// Vectorize inserts a string value as a searchable vector with the specified weight.
func Vectorize(weight string) string {
	return fmt.Sprintf("setweight(to_tsvector('english', unaccent(?)), '%s')", weight)
}

type WithBeforeUpdate interface {
	BeforeUpdate(query *bun.UpdateQuery) *bun.UpdateQuery
}

type WithBeforeCreate interface {
	BeforeCreate(query *bun.InsertQuery) *bun.InsertQuery
}
