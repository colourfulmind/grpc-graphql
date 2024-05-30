package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	conn := "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("could not create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("could not create migrate instance: %v", err)
	}

	var flg string
	flag.StringVar(&flg, "migrate", "", "migrate command")
	flag.Parse()

	switch flg {
	case "up":
		if err = m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("could not run up migrations: %v", err)
		}
	case "down":
		if err = m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("could not run down migrations: %v", err)
		}
	default:
		log.Fatalf("unknown migrate command: %v", flg)
	}

	log.Println("migrations applied successfully")
}
