package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"server/pkg/core/actions"
	db "server/pkg/core/storage"

	zlog "github.com/rs/zerolog/log"
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

func LoginRequestFunction(storage *db.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyRaw, err := ioutil.ReadAll(r.Body)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, LoginResponse{Message: err.Error()}, err, "")
			return
		}

		zlog.Debug().Bytes("request_body", bodyRaw).Str("method", r.Method).Send()

		if r.Method == http.MethodOptions {
			SendResponse(w, http.StatusOK, nil, nil, "CORS Request processing.")
			return
		}

		creds := LoginCreds{}
		err = json.Unmarshal(bodyRaw, &creds)
		if err != nil {
			SendResponse(w, http.StatusBadRequest, LoginResponse{Message: err.Error()}, err, "")
			return
		}

		foundUser, err := actions.GetUser(storage, creds.Login)
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
			token, err := actions.IssueToken(storage, creds.Login)
			if err != nil {
				SendResponse(w, http.StatusInternalServerError, nil, err, "Error during token issue.")
			}
	
			http.SetCookie(w, &http.Cookie{
				Name:  "token",
				Value: token,
			})
			
			SendResponse(w, http.StatusOK, nil, nil, "")
		}
	}

}
