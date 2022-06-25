package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	// CORS Request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	creds := LoginCreds{}
	bodyRaw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(bodyRaw))

	err = json.Unmarshal(bodyRaw, &creds)
	if err != nil {
		log.Fatal(err)
	}

	responseBody, err := json.Marshal(&LoginResponse{
		Message: "Not implemented yet",
	})
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set(contentTypeHeader, jsonContentType)
	w.WriteHeader(http.StatusNotImplemented)
	w.Write(responseBody)
}
