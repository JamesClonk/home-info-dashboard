package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

type PostgresAdapter struct {
	Database *sql.DB
	URI      string
	Type     string
}

func newPostgresAdapter(uri string) *PostgresAdapter {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}
	return &PostgresAdapter{
		Database: db,
		URI:      uri,
		Type:     "postgres",
	}
}

func (adapter *PostgresAdapter) GetDatabase() *sql.DB {
	return adapter.Database
}

func (adapter *PostgresAdapter) GetURI() string {
	return adapter.URI
}

func (adapter *PostgresAdapter) GetType() string {
	return adapter.Type
}

func (adapter *PostgresAdapter) RunMigrations(basePath string) error {
	driver, err := postgres.WithInstance(adapter.Database, &postgres.Config{})
	if err != nil {
		log.Println("Could not create database migration driver")
		log.Fatalln(err)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s/postgres", basePath), "postgres", driver)
	if err != nil {
		log.Println("Could not create database migration instance")
		log.Fatalln(err)
	}

	return m.Up()
}
