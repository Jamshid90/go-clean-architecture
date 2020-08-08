package main

import (
	"context"
	"flag"
	"github.com/Jamshid90/go-clean-architecture/pkg/auth"
	"github.com/Jamshid90/go-clean-architecture/pkg/config"
	"github.com/Jamshid90/go-clean-architecture/pkg/http/rest/middleware"
	"github.com/Jamshid90/go-clean-architecture/pkg/http/server"
	zaplogger "github.com/Jamshid90/go-clean-architecture/pkg/logger"
	"github.com/Jamshid90/go-clean-architecture/pkg/refreshtoken"
	"github.com/Jamshid90/go-clean-architecture/pkg/user"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

var (
	configPath = flag.String("cp", "./config.toml", "path to configuration file")
	logLevel   = flag.String("log-level", "debug", "allowed value for logger level: debug, info, warn, error, error, dpanic, panic, fatal")
)

func main() {

	flag.Parse()

	// initialization config
	config, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// connect  pgxpool
	dbpool, err := pgxpool.Connect(context.Background(), config.GetPsqlConnStr())
	if err != nil {
		log.Fatal(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	// Initialization zap logger
	logger, err := zaplogger.NewDevZapLogger(*logLevel)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	r := chi.NewRouter()

	// initialization repositorys
	userRepo := user.NewPgxUserRepository(dbpool)
	refreshTokenRepo := refreshtoken.NewRefreshTokenRepositoryPgx(dbpool)

	// initialization usecase
	userUsecase := user.NewUserUsecase(userRepo, config.Context.Timeout)
	refreshTokenUsecase := refreshtoken.NewRefreshTokenUsecase(refreshTokenRepo, config.Context.Timeout)

	r.Route("/api", func(r chi.Router) {

		// initialization api middleware
		r.Use(middleware.RequestID)
		r.Use(middleware.Cors)
		r.Use(middleware.ContentTypeJson)
		r.Use(middleware.Logger(logger))

		// initialization auth handlers
		auth.NewAuthHandler(r, &userUsecase, &refreshTokenUsecase, config, logger)

		// initialization user handlers
		user.NewUserHandler(r, &userUsecase, logger)

	})

	// initialization server
	appServer := server.NewServer(config, r)
	logger.Info("Listen: http://" + appServer.GetServerAddr())
	log.Fatal(appServer.Run())
}
