package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func newMysqlAdapter(uri string) *Adapter {
	db, err := sql.Open("mysql", uri)
	if err != nil {
		panic(err)
	}
	return &Adapter{
		Database: db,
		URI:      uri,
		Type:     "mysql",
	}
}
