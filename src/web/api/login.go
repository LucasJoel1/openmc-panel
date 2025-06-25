package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"openmc-panel/src/db"
	"strconv"
	"time"
)

type LoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	Username    string `json:"username"`
	Permissions string `json:"permissions"`
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var reqBody LoginRequestBody

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request format"))
		return
	}

	token, permissions, err := db.VerifyUserCreds(reqBody.Username, reqBody.Password)

	if errors.Is(err, db.ErrUserNotExist) || errors.Is(err, db.ErrWrongPassword) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid username or password"))
		return
	}

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("an unknown server error has occured"))
		log.Println(err)
		return
	}

	userInfo := UserInfo{
		Username:    reqBody.Username,
		Permissions: strconv.FormatUint(permissions, 10),
	}

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userInfo)
}
