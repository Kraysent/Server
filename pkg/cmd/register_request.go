package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"server/pkg/core/actions"
	"server/pkg/db"
)

type RegisterCreds struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

func RegisterRequestFunction(storage *db.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		action := actions.NewStorageAction(storage)

		_, err = action.CreateUser(creds.Login, creds.Password, creds.Description)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, RegisterResponse{Message: err.Error()}, err, "")
			return
		}

		token, err := action.IssueToken(creds.Login)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, nil, err, "Error during token issue.")
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: token,
		})
		SendResponse(w, http.StatusOK, nil, nil, "")
	})
}
