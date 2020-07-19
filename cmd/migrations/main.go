package main

import (
	"os"
	"fmt"
	"log"
	"flag"
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/Jamshid90/go-clean-architecture/pkg/config"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

var (
	configPath = flag.String("config-path", "./config.toml", "path to configuration file")
	action = flag.String("action", "up", "available action: up, down, drop")
)

func main() {
	flag.Parse()

	// initialization config
	config, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	//connect database
	db, err := sql.Open("pgx", config.GetPsqlConnStr())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "./db/migrations"),
		config.Database.DBName, driver)

	switch *action {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "drop":
		err = m.Drop()
	}

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Success")
}
