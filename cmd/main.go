package main

import (
	"fmt"
	"github.com/dewi911/cruda-app/internal/config"
	"github.com/dewi911/cruda-app/internal/repository/psql"
	"github.com/dewi911/cruda-app/internal/service"
	"github.com/dewi911/cruda-app/internal/transport/grpc"
	"github.com/dewi911/cruda-app/internal/transport/rest"
	"github.com/dewi911/cruda-app/pkg/database"
	"github.com/dewi911/cruda-app/pkg/hash"
	_ "github.com/lib/pq"
	"os"

	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	//var crftest = new(config.Config)
	log.Printf("config: %+v\n", cfg.DB)

	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.DB.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//init deps
	hasher := hash.NewSHA1Hasher("salt")

	bookRepo := psql.NewBooks(db)
	bookService := service.NewBooks(bookRepo)

	usersRepo := psql.NewUsers(db)
	tokensRepo := psql.NewTokens(db)

	auditClient, err := grpc.NewClient(9000)
	if err != nil {
		log.Fatal(err)
	}

	usersService := service.NewUsers(usersRepo, tokensRepo, auditClient, hasher, []byte("sample secret key"), cfg.Auth.TokenTTL)

	handler := rest.NewHandler(bookService, usersService)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: handler.InitRouter(),
	}

	log.Info("Listening on port 8080")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
