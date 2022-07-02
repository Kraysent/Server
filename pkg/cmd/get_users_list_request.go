package cmd

import (
	"net/http"
	"server/pkg/core/actions"
	db "server/pkg/core/storage"

	zlog "github.com/rs/zerolog/log"
)

func GetUserByLoginRequestFunction(storage *db.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login := r.URL.Query().Get("login")

		action := actions.NewStorageAction(storage)
		user, err := action.GetUser(login)
		if err != nil {
			zlog.Error().Err(err).Send()
		}

		SendResponse(w, http.StatusOK, user, nil, "")
	})
}
