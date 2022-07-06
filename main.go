package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"server/pkg/cmd"
	"server/pkg/core"
	"server/pkg/core/server"
	"server/pkg/core/server/middleware"
	"server/pkg/db"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	config, err := core.NewConfig("configs/dev.yaml")
	if err != nil {
		zlog.Fatal().Err(err).Msg("Error during config loading.")
	}

	level, err := zerolog.ParseLevel(config.Server.Level)
	if err != nil {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	storage := db.NewStorage(config.Storage)
	err = storage.Connect()
	if err != nil {
		zlog.Fatal().Err(err).Msg("Error occured during connection to storage")
	}
	defer storage.Disconnect()

	router := server.NewRouter()
	router.AddMiddleware(middleware.CORSMiddleware)
	router.AddMiddleware(middleware.LoggingMiddleware)

	// Common handlers
	router.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "pong") }))
	router.Handle("/login", cmd.LoginRequestFunction(storage))
	router.Handle("/register", cmd.RegisterRequestFunction(storage))
	router.Handle("/profile", cmd.ProfileRequestFunction(storage), middleware.GetAuthMiddleware(storage))
	router.Handle("/search", cmd.SearchRequestFunction(storage), middleware.GetAuthMiddleware(storage))

	// Admin handlers
	router.Handle("/get_user", cmd.GetUserByLoginRequestFunction(storage))

	zlog.Info().Str("address", fmt.Sprintf("http://127.0.0.1:%d", config.Server.Port)).Msg("Server is listening")
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), nil)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Error occured during listening")
	}
}
