package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"server/pkg/actions"
	"server/pkg/storage"
)

const (
	contentTypeHeader = "Content-Type"

	jsonContentType = "application/json"
)

type LoginCreds struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
}

func LoginRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	w.Header().Set(contentTypeHeader, jsonContentType)

	// CORS Request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	bodyRaw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(RegisterResponse{Message: err.Error()})
		w.Write(resp)
		log.Println(err)
		return
	}

	creds := LoginCreds{}
	err = json.Unmarshal(bodyRaw, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(RegisterResponse{Message: err.Error()})
		w.Write(resp)
		log.Printf("%s: %s", err, string(bodyRaw))
		return
	}

	foundUser, err := storage.Get(creds.Login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(RegisterResponse{Message: fmt.Sprintf("User not found: %s", err.Error())})
		w.Write(resp)
		log.Printf("%s: %s", err, string(bodyRaw))
		return
	}

	if foundUser.PasswordHash != actions.HashPassword(creds.Password, foundUser.Salt) {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
