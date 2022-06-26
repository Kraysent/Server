package cmd

import (
	"encoding/json"
	"log"
	"net/http"
	"server/pkg/storage"
)

func GetUserByLoginRequest(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query().Get("login")

	user, err := storage.Get(login)
	if err != nil {
		log.Fatal(err)
	}

	body, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(body)
}
