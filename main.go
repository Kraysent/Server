package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"server/pkg/cmd"
	"server/pkg/core"
	db "server/pkg/core/storage"
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

	// Common handlers
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "pong") })
	http.HandleFunc("/login", cmd.LoginRequestFunction(storage))
	http.HandleFunc("/register", cmd.RegisterRequestFunction(storage))
	http.HandleFunc("/profile", cmd.ProfileRequestFunction(storage))

	// Admin handlers
	http.HandleFunc("/get_user", cmd.GetUserByLoginRequestFunction(storage))

	zlog.Info().Str("address", fmt.Sprintf("http://127.0.0.1:%d", config.Server.Port)).Msg("Server is listening")
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), nil)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Error occured during listening")
	}
}
