package cmd

import (
	"net/http"
	"server/pkg/core/actions"
	"server/pkg/db"
	"time"
)

func ProfileRequestFunction(storage *db.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			SendResponse(w, http.StatusUnauthorized, nil, err, "")
			return
		}

		action := actions.NewStorageAction(storage)

		login := r.URL.Query().Get("login")

		tokenValid, err := action.CheckUserToken(login, tokenCookie.Value, time.Now())
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, nil, err, "Error during token check.")
			return
		}
		if !tokenValid {
			SendResponse(w, http.StatusUnauthorized, nil, nil, "")
			return
		}

		user, err := action.GetUser(login)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, nil, err, "Error during user query.")
			return
		}

		SendResponse(w, http.StatusOK, user, nil, "")
	})
}
