package cmd

import (
	"net/http"
	"server/pkg/core/actions"
	db "server/pkg/core/storage"

	zlog "github.com/rs/zerolog/log"
)

func GetUserByLoginRequestFunction(storage *db.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		login := r.URL.Query().Get("login")

		user, err := actions.GetUser(storage, login)
		if err != nil {
			zlog.Error().Err(err).Send()
		}

		SendResponse(w, http.StatusOK, user, nil, "")
	}
}
