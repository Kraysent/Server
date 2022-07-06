package middleware

import (
	"net/http"
	"server/pkg/core/actions"
	"server/pkg/db"
	"time"
)

func GetAuthMiddleware(storage *db.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenCookie, err := r.Cookie("token")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			loginCookie, err := r.Cookie("login")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			action := actions.NewStorageAction(storage)

			tokenValid, err := action.CheckUserToken(loginCookie.Value, tokenCookie.Value, time.Now())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !tokenValid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
