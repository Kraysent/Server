package cmd

import (
	zlog "github.com/rs/zerolog/log"
	"net/http"
	"server/pkg/storage"
)

func GetUserByLoginRequest(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query().Get("login")

	user, err := storage.Get(login)
	if err != nil {
		zlog.Error().Err(err).Send()
	}

	SendResponse(w, http.StatusOK, user, nil, "")
}
