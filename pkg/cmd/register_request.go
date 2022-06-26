package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"server/pkg/core/actions"
	db "server/pkg/core/storage"
)

type RegisterCreds struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

func RegisterRequestFunction(storage *db.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			SendResponse(w, http.StatusOK, nil, nil, "CORS Request processing.")
			return
		}

		bodyRaw, err := ioutil.ReadAll(r.Body)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, RegisterResponse{Message: err.Error()}, err, "")
			return
		}

		creds := RegisterCreds{}
		err = json.Unmarshal(bodyRaw, &creds)
		if err != nil {
			SendResponse(w, http.StatusBadRequest, RegisterResponse{Message: err.Error()}, err, "")
			return
		}

		_, err = actions.CreateUser(storage, creds.Login, creds.Password, creds.Description)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, RegisterResponse{Message: err.Error()}, err, "")
			return
		}

		SendResponse(w, http.StatusOK, nil, nil, "")
	}
}
