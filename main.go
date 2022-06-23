package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	port         = 8081
	user         = "testserver"
	password     = "passw0rd"
	postgresPort = 5432
	database     = "serverdb"
)

func main() {
	// Common handlers
	http.HandleFunc("/", RootRequest)
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "pong") })
	http.HandleFunc("/login", LoginRequest)

	// Admin handlers
	http.HandleFunc("/add_test_user", AddUserRequestFunction(user, password, postgresPort, database))
	http.HandleFunc("/get_users_list", GetUsersRequestFunction(user, password, postgresPort, database))

	log.Printf("Listening on http://127.0.0.1:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}