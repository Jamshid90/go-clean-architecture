package main

import (
	"context"
	"flag"
	"github.com/Jamshid90/go-clean-architecture/pkg/config"
	"github.com/Jamshid90/go-clean-architecture/pkg/middleware"
	"github.com/Jamshid90/go-clean-architecture/pkg/server"
	"github.com/Jamshid90/go-clean-architecture/pkg/user"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var (
	configPath = flag.String("cp", "./config.toml", "path to configuration file")
)

func main() {

	flag.Parse()

	// initialization config
	config, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	// connect pgx
	conn, err := pgx.Connect(context.Background(), config.GetPsqlConnStr())
	if err != nil {
		log.Fatal(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	/*
	//connect pq
	conn, err := sql.Open("postgres", config.GetPsqlConnStr())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	*/

	r := chi.NewRouter()

	// initialization repositorys
	userRepo := user.NewPgxUserRepository(conn)

	// initialization usecase
	userUsecase := user.NewUserUsecase(userRepo, config.Context.Timeout)

	r.Route("/api", func(r chi.Router) {

		// initialization api middleware
		r.Use(middleware.Cors)
		r.Use(middleware.ContentTypeJson)

		// initialization api handlers
		user.NewUserHandler(r, &userUsecase)
	})

	// initialization server
	appServer := server.NewServer(config, r)
	log.Println("Listen:", "http://"+appServer.GetServerAddr())
	log.Fatal(appServer.Run())
}