package main

import (
	"fmt"
	"net/http"
	"os"
	"server/pkg/cmd"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

const (
	port = 8081
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Common handlers
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "pong") })
	http.HandleFunc("/login", cmd.LoginRequest)
	http.HandleFunc("/register", cmd.RegisterRequest)

	// Admin handlers
	http.HandleFunc("/get_user", cmd.GetUserByLoginRequest)

	zlog.Info().Str("address", fmt.Sprintf("http://127.0.0.1:%d", port)).Msg("Server is listening")
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Error occured during listening")
	}
}
