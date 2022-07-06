package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"server/pkg/core/actions"
	"server/pkg/db"
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

func LoginRequestFunction(storage *db.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyRaw, err := ioutil.ReadAll(r.Body)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, LoginResponse{Message: err.Error()}, err, "")
			return
		}
		creds := LoginCreds{}
		err = json.Unmarshal(bodyRaw, &creds)
		if err != nil {
			SendResponse(w, http.StatusBadRequest, LoginResponse{Message: err.Error()}, err, "")
			return
		}

		action := actions.NewStorageAction(storage)

		foundUser, err := action.GetUser(creds.Login)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, LoginResponse{Message: err.Error()}, err, "")
			return
		}
		if foundUser == nil {
			SendResponse(w, http.StatusUnauthorized, nil, nil, "")
			return
		}

		if foundUser.PasswordHash != actions.HashPassword(creds.Password, foundUser.Salt) {
			SendResponse(w, http.StatusUnauthorized, nil, nil, "")
		} else {
			token, err := action.IssueToken(creds.Login)
			if err != nil {
				SendResponse(w, http.StatusInternalServerError, nil, err, "Error during token issue.")
			}

			http.SetCookie(w, &http.Cookie{
				Name:  "token",
				Value: token,
			})

			SendResponse(w, http.StatusOK, nil, nil, "")
		}
	})
}
