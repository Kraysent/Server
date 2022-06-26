package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"server/pkg/actions"
)

type RegisterCreds struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

func RegisterRequest(w http.ResponseWriter, r *http.Request) {
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

	creds := RegisterCreds{}
	err = json.Unmarshal(bodyRaw, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(RegisterResponse{Message: err.Error()})
		w.Write(resp)
		log.Printf("%s: %s", err, string(bodyRaw))
		return
	}

	_, err = actions.CreateUser(creds.Login, creds.Password, creds.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(RegisterResponse{Message: err.Error()})
		w.Write(resp)
		log.Printf("%s: %s", err, string(bodyRaw))
		return
	}

	w.WriteHeader(http.StatusOK)
}
