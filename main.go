package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	port = 8081
)

func main() {
	// Common handlers
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "pong") })
	http.HandleFunc("/login", LoginRequest)
	http.HandleFunc("/register", RegisterRequest)

	// Admin handlers
	http.HandleFunc("/get_user", GetUserByLoginRequest)

	log.Printf("Listening on http://127.0.0.1:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
