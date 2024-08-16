package app

import (
	"auth/internal/db"
	"auth/internal/handlers"
	"auth/internal/repo"
	"auth/internal/services"
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

func StartApp() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	db, err := db.DBinit(ctx)
	if err != nil {
		log.Fatal(err)
	}
	repository := repo.NewRepo(db)

	services := services.NewService(repository)

	handlers := handlers.NewHandler(services)

	finalHandler := http.HandlerFunc(handlers.Final)

	http.HandleFunc("/login", handlers.GetTokenHandler)
	http.HandleFunc("/refresh", handlers.RefreshTokenHandler)
	http.Handle("/check", handlers.MiddlewareOne(finalHandler))

	err = http.ListenAndServe(os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
