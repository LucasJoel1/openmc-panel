package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"openmc-panel/src/db"
)

type RegistrationRequestBody struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var reqBody RegistrationRequestBody

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request format"))
		return
	}

	if reqBody.Password != reqBody.ConfirmPassword {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("passwords do not match"))
		return
	}

	err = db.CreateUser(reqBody.Username, reqBody.Password, false)

	if errors.Is(err, db.ErrUsernameInvalid) || errors.Is(err, db.ErrPasswordInvalid) || errors.Is(err, db.ErrUserExists) { 
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("an unknown server error has occured"))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}