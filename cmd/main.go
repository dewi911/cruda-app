package main

import (
	"cruda-app/internal/repository/psql"
	"cruda-app/internal/service"
	"cruda-app/internal/transport/rest"
	"cruda-app/pkg/database"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

func main() {
	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
		Password: "postgres",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bookRepo := psql.NewBooks(db)
	bookService := service.NewBooks(bookRepo)
	handler := rest.NewHandler(bookService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.InitRouter(),
	}

	log.Println("Listening on port 8080", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	}
}
