package cmd

import (
	"net/http"
	"server/pkg/core/actions"
	db "server/pkg/core/storage"

	zlog "github.com/rs/zerolog/log"
)

func ProfileRequestFunction(storage *db.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		zlog.Debug().Str("method", r.Method).Send()

		if r.Method == http.MethodOptions {
			SendResponse(w, http.StatusOK, nil, nil, "CORS Request processing.")
			return
		}

		tokenCookie, err := r.Cookie("token")
		if err != nil {
			SendResponse(w, http.StatusUnauthorized, nil, err, "")
			return
		}

		login := r.URL.Query().Get("login")

		tokenValid, err := actions.CheckUserToken(storage, login, tokenCookie.Value)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, nil, err, "Error during token check.")
			return
		}
		if !tokenValid {
			SendResponse(w, http.StatusUnauthorized, nil, nil, "")
			return
		}

		user, err := actions.GetUser(storage, login)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, nil, err, "Error during user query.")
			return
		}

		SendResponse(w, http.StatusOK, user, nil, "")
	}
}
