package api

import (
	"net/http"
	"openmc-panel/src/db"
	"strconv"
)

func GetUserPermissions(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if username == "" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("username empty"))
		return
	}

	perms, err := db.GetUserPermissions(username)

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error retrieving user perms"))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatUint(perms, 10)))
}