package main

import (
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/user/urlshortener/internal/core"
	"github.com/user/urlshortener/internal/persistence"
	"github.com/user/urlshortener/internal/server"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	db, err := persistence.NewSQLitePersistence("urls.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize persistence")
	}

	core := core.NewCore(db)
	server := server.NewServer(core)

	router := httprouter.New()
	server.RegisterRoutes(router)

	log.Info().Msg("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
