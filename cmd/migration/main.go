package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/Jamshid90/go-clean-architecture/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
)

var (
	configPath = flag.String("config-path", "./config.toml", "path to configuration file")
	action     = flag.String("action", "up", "available action: up, down, drop")
)

func main() {
	flag.Parse()

	// initialization config
	config, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatal("config", err)
	}

	//connect database
	db, err := sql.Open("pgx", config.GetPsqlConnStr())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("ping", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Println("driver", err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "./db/migrations"),
		config.Database.DBName, driver)
	if err != nil {
		log.Println("new with database instance", err)
		return
	}

	switch *action {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "drop":
		err = m.Drop()
	}

	if err != nil {
		log.Println("action", err)
		return
	}

	fmt.Println("Success")
}
