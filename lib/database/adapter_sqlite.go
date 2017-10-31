package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/mattes/migrate"
	sqlite "github.com/mattes/migrate/database/sqlite3"
	_ "github.com/mattes/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteAdapter struct {
	Database *sql.DB
	URI      string
	Type     string
}

func newSQLiteAdapter(uri string) *SQLiteAdapter {
	filename := strings.SplitN(uri, "sqlite3://", 2)
	if len(filename) != 2 {
		panic("invalid sqlite3:// uri")
	}

	db, err := sql.Open("sqlite3", filename[1])
	if err != nil {
		panic(err)
	}
	return &SQLiteAdapter{
		Database: db,
		URI:      uri,
		Type:     "sqlite",
	}
}

func (adapter *SQLiteAdapter) GetDatabase() *sql.DB {
	return adapter.Database
}

func (adapter *SQLiteAdapter) GetURI() string {
	return adapter.URI
}

func (adapter *SQLiteAdapter) GetType() string {
	return adapter.Type
}

func (adapter *SQLiteAdapter) RunMigrations(basePath string) error {
	driver, err := sqlite.WithInstance(adapter.Database, &sqlite.Config{})
	if err != nil {
		log.Println("Could not create database migration driver")
		log.Fatalln(err)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s/sqlite", basePath), "sqlite3", driver)
	if err != nil {
		log.Println("Could not create database migration instance")
		log.Fatalln(err)
	}

	return m.Up()
}
