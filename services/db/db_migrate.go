package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var postgresURI string

func MigrateDatabase(pgURI, migrationPath string, migrationVersion uint) error {
	postgresURI = fmt.Sprintf("%s?connect_timeout=10&sslmode=disable", pgURI)
	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		log.Printf("Error while opening connection: %s\n", err.Error())
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Printf("Error creating db driver: %s\n", err.Error())
		return err
	}

	if migrationPath == "" {
		migrationPath = "file://db/migrations/"
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", driver)
	if err != nil {
		log.Printf("Error creating migrate instance: %s\n", err)
		return err
	}
	defer m.Close()

	//Give the version number which we want to migrate
	err = m.Migrate(migrationVersion)
	if err != nil && err != migrate.ErrNoChange {
		log.Printf("Error migrating database: %s\n", err.Error())
		return err
	}
	return nil
}
