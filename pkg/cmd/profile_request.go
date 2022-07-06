package cmd

import (
	"net/http"
	"server/pkg/core/actions"
	"server/pkg/db"
)

func ProfileRequestFunction(storage *db.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		action := actions.NewStorageAction(storage)

		loginCookie, err := r.Cookie("login")
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, nil, err, "Error during cookie reading.")
		}

		user, err := action.GetUser(loginCookie.Value)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, nil, err, "Error during user query.")
			return
		}

		SendResponse(w, http.StatusOK, user, nil, "")
	})
}
