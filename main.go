package main

import (
	"fmt"
	"net/http"
	"os"
	"server/pkg/cmd"
	db "server/pkg/core/storage"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

const (
	serverPort = 8081
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	storage := db.NewStorage(db.StorageConfig{
		User:     "testserver",
		Password: "passw0rd",
		Port:     5432,
		DBName:   "serverdb",
	})

	// Common handlers
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "pong") })
	http.HandleFunc("/login", cmd.LoginRequestFunction(storage))
	http.HandleFunc("/register", cmd.RegisterRequestFunction(storage))

	// Admin handlers
	http.HandleFunc("/get_user", cmd.GetUserByLoginRequestFunction(storage))

	zlog.Info().Str("address", fmt.Sprintf("http://127.0.0.1:%d", serverPort)).Msg("Server is listening")
	err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Error occured during listening")
	}
}
