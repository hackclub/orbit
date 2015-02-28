package datastore

import (
	"log"
	"os"
	"sync"

	"github.com/jmoiron/modl"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB = &modl.DbMap{Dialect: modl.PostgresDialect{}}

var DBH modl.SqlExecutor = DB

var connectOnce sync.Once

func Connect() {
	connectOnce.Do(func() {
		var err error
		DB.Dbx, err = sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatal("Error connecting to PostgreSQL database (using DATABASE_URL environment variable): ", err)
		}
		DB.Db = DB.Dbx.DB
	})
}

var createSQL []string

// Create the database schema. It calls log.Fal if it encounters an error.
func Create() {
	if err := DB.CreateTablesIfNotExists(); err != nil {
		log.Fatal("Error creating tables: ", err)
	}
	for _, query := range createSQL {
		if _, err := DB.Exec(query); err != nil {
			log.Fatalf("Error running query %q: %s", query, err)
		}
	}
}

// Drop the database schema.
func Drop() {
	DB.DropTables()
}

// transact calls fn in a DB transaction. If dbh is a transaction, then it just
// calls the function. Otherwise, it begins a transaction, rolling back on
// failure and committing on success.
func transact(dbh modl.SqlExecutor, fn func(dbh modl.SqlExecutor) error) error {
	var sharedTx bool
	tx, sharedTx := dbh.(*modl.Transaction)
	if !sharedTx {
		var err error
		tx, err = dbh.(*modl.DbMap).Begin()
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
	}

	if err := fn(tx); err != nil {
		return err
	}

	if !sharedTx {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
