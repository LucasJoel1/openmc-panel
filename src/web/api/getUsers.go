package api

import (
	"encoding/json"
	"net/http"
	"openmc-panel/src/db"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetUsersNames()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("an unknown server error has occured"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}