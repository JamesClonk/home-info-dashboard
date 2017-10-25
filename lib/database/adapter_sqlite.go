package database

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func newSQLiteAdapter(uri string) *Adapter {
	filename := strings.SplitN(uri, "sqlite3://", 2)
	if len(filename) != 2 {
		panic("invalid sqlite3:// uri")
	}

	db, err := sql.Open("sqlite3", filename[1])
	if err != nil {
		panic(err)
	}
	return &Adapter{
		Database: db,
		URI:      uri,
		Type:     "sqlite",
	}
}
