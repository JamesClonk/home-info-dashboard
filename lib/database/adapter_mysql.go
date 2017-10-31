package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/mysql"
	_ "github.com/mattes/migrate/source/file"
)

type MysqlAdapter struct {
	Database *sql.DB
	URI      string
	Type     string
}

func newMysqlAdapter(uri string) *MysqlAdapter {
	db, err := sql.Open("mysql", uri)
	if err != nil {
		panic(err)
	}
	return &MysqlAdapter{
		Database: db,
		URI:      uri,
		Type:     "mysql",
	}
}

func (adapter *MysqlAdapter) GetDatabase() *sql.DB {
	return adapter.Database
}

func (adapter *MysqlAdapter) GetURI() string {
	return adapter.URI
}

func (adapter *MysqlAdapter) GetType() string {
	return adapter.Type
}

func (adapter *MysqlAdapter) RunMigrations(basePath string) error {
	driver, err := mysql.WithInstance(adapter.Database, &mysql.Config{})
	if err != nil {
		log.Println("Could not create database migration driver")
		log.Fatalln(err)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s/mysql", basePath), "mysql", driver)
	if err != nil {
		log.Println("Could not create database migration instance")
		log.Fatalln(err)
	}

	return m.Up()
}
