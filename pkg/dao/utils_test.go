package dao_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/in-rich/uservice-search/migrations"
	_ "github.com/in-rich/uservice-search/migrations"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"time"
)

func OpenDB() *bun.DB {
	dsn := "postgres://test:test@localhost:5432/test?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	err := db.Ping()
	for i := 0; i < 10 && err != nil; i++ {
		time.Sleep(1 * time.Second)
		err = db.Ping()
	}

	// Just in case something went wrong on latest run.
	ClearDB(db)

	if err := migrations.Migrate(db); err != nil {
		panic(err)
	}

	return db
}

func ClearDB(db *bun.DB) {
	if _, err := db.ExecContext(context.TODO(), "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"); err != nil {
		panic(err)
	}

	if _, err := db.ExecContext(context.TODO(), "GRANT ALL ON SCHEMA public TO public;"); err != nil {
		panic(err)
	}
	if _, err := db.ExecContext(context.TODO(), "GRANT ALL ON SCHEMA public TO test;"); err != nil {
		panic(err)
	}
}

func CloseDB(db *bun.DB) {
	ClearDB(db)

	if err := db.Close(); err != nil {
		panic(err)
	}
}

func BeginTX[T any](db bun.IDB, fixtures []T) bun.Tx {
	tx, err := db.BeginTx(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	for _, fixture := range fixtures {
		var query *bun.InsertQuery

		if creator, ok := interface{}(fixture).(entities.WithBeforeCreate); ok {
			query = creator.BeforeCreate(tx.NewInsert().Model(creator))
		} else {
			query = tx.NewInsert().Model(fixture)
		}

		if _, err := query.Exec(context.TODO()); err != nil {
			panic(fmt.Errorf("failed to insert fixture: %w\n%s\n", err, query))
		}
	}

	return tx
}

func RollbackTX(tx bun.Tx) {
	_ = tx.Rollback()
}
