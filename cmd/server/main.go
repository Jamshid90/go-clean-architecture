package main

import (
	"log"
	"flag"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/go-chi/chi"
	"github.com/Jamshid90/go-clean-architecture/pkg/user"
	"github.com/Jamshid90/go-clean-architecture/pkg/config"
	"github.com/Jamshid90/go-clean-architecture/pkg/server"
	"github.com/Jamshid90/go-clean-architecture/pkg/middleware"
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

	// initialization database
	db, err := sql.Open("postgres", config.GetPsqlConnStr())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	// initialization repositorys
	//resourceRepo := postgresql.NewPostgresqlResourceRepository(db)
	userRepo     := user.NewPostgresqlUserRepository(db)

	// initialization usecase
	//rUsecase     := usecase.NewResourceUsecase(resourceRepo)
	userUsecase  := user.NewUserUsecase(userRepo)

	r.Route("/api", func(r chi.Router) {
		// initialization api handlers
		//http.NewresourceHandler(r, &rUsecase)
		r.Use(middleware.Cors)
		r.Use(middleware.ContentTypeJson)

		user.NewUserHandler(r, &userUsecase)
	})

	// initialization server
	appServer := server.NewServer(config, r)
	log.Println("Listen:", "http://"+appServer.GetServerAddr())
	log.Fatal(appServer.Run())
}