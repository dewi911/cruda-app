package main

import (
	"cruda-app/internal/config"
	"cruda-app/internal/repository/psql"
	"cruda-app/internal/service"
	"cruda-app/internal/transport/rest"
	"cruda-app/pkg/database"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

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

	bookRepo := psql.NewBooks(db)
	bookService := service.NewBooks(bookRepo)
	handler := rest.NewHandler(bookService)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: handler.InitRouter(),
	}

	log.Println("Listening on port 8080", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
