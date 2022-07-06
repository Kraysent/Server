package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"server/pkg/core/actions"
	"server/pkg/db"
)

type SearchRequest struct {
	LoginPattern string `json:"pattern"`
}

type User struct {
	Login       string `json:"login"`
	Description string `json:"description"`
}

type SearchResponse struct {
	Users []User `json:"users"`
}

func SearchRequestFunction(storage *db.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		action := actions.NewStorageAction(storage)

		bodyRaw, err := ioutil.ReadAll(r.Body)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, RegisterResponse{Message: err.Error()}, err, "")
			return
		}

		request := SearchRequest{}
		err = json.Unmarshal(bodyRaw, &request)
		if err != nil {
			SendResponse(w, http.StatusBadRequest, RegisterResponse{Message: err.Error()}, err, "")
			return
		}

		users, err := action.FindUserByPattern(request.LoginPattern)
		if err != nil {
			SendResponse(w, http.StatusInternalServerError, nil, err, "Error during user query.")
			return
		}

		response := SearchResponse{
			Users: []User{},
		}
		for _, user := range users {
			response.Users = append(response.Users, User{
				Login: user.Login, Description: user.Description,
			})
		}

		SendResponse(w, http.StatusOK, response, nil, "")
	})
}
