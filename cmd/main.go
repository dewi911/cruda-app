package main

import (
	"cruda-app/internal/config"
	"cruda-app/internal/repository/psql"
	"cruda-app/internal/service"
	"cruda-app/internal/transport/rest"
	"cruda-app/pkg/database"
	"cruda-app/pkg/hash"
	"fmt"
	"github.com/kelseyhightower/envconfig"
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

	var crftest = new(config.Config)

	if err := envconfig.Process("db", &crftest.DB); err != nil {

		log.Printf("Error processing DB_HOST: %v", err)

	}
	fmt.Printf("%+v\n", crftest.DB.Host)

	log.Printf("config: %+v\n", crftest.DB)

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
	usersService := service.NewUsers(usersRepo, hasher)

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
