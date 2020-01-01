package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/JamesClonk/home-info-dashboard/lib/env"
	cfenv "github.com/cloudfoundry-community/go-cfenv"
)

type Adapter interface {
	GetDatabase() *sql.DB
	GetURI() string
	GetType() string
	RunMigrations(string) error
}

func NewAdapter() (db Adapter) {
	var databaseType, databaseUri string

	// get db type
	databaseType = env.Get("HOME_INFO_DB_TYPE", "postgres")

	// check for VCAP_SERVICES first
	vcap, err := cfenv.Current()
	if err != nil {
		log.Println("Could not parse VCAP environment variables")
		log.Println(err)
	} else {
		service, err := vcap.Services.WithName("home_info_db")
		if err != nil {
			log.Println("Could not find home_info_db service in VCAP_SERVICES")
			log.Fatal(err)
		}
		databaseUri = fmt.Sprintf("%v", service.Credentials["uri"])
	}

	// if HOME_INFO_DB_URI is not yet set then try to read it from ENV
	if len(databaseUri) == 0 {
		databaseUri = env.MustGet("HOME_INFO_DB_URI")
	}

	// setup database adapter
	switch databaseType {
	case "postgres":
		db = newPostgresAdapter(databaseUri)
	case "sqlite":
		db = newSQLiteAdapter(databaseUri)
	default:
		log.Fatalf("Invalid database type: %s\n", databaseType)
	}

	// panic if no database adapter was set up
	if db == nil {
		log.Fatal("Could not set up database adapter")
	}
	return db
}
