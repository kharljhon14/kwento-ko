package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"www.github.com/kharljhon14/kwento-ko/cmd/api"
)

func main() {
	_, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	clientCallbackURL := os.Getenv("GOOGLE_CALLBACK_URL")

	if clientID == "" || clientSecret == "" || clientCallbackURL == "" {
		log.Fatal("env (GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GOOGLE_CALLBACK_URL) are required")
	}

	goth.UseProviders(
		google.New(clientID, clientSecret, clientCallbackURL),
	)

	server, err := api.NewServer()
	if err != nil {
		log.Fatal("failed to create new server:", err)
	}

	err = server.Start(os.Getenv("SERVER_ADDRESS"))
	if err != nil {
		log.Fatal("failed to start server:", err)
	}
}
