package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"openmc-panel/src/db"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

type ModifyUserPermissionsBody struct {
	Username    string `json:"username"`
	Permissions string `json:"permissions"`
}

func ModifyUserPermissions(w http.ResponseWriter, r *http.Request) {
	var reqBody ModifyUserPermissionsBody

	cookie, err := r.Cookie("auth_token")

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("token not found, user unauthorized"))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request format"))
		return
	}

	var secret []byte = db.GetJWTSecret()

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid token"))
		return
	}

	claims, validToken := token.Claims.(jwt.MapClaims)

	if !validToken {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid token"))
		return
	}

	// Handle permissions claim as float64 (JWT standard) or string
	var perms uint64
	switch v := claims["permissions"].(type) {
	case string:
		var err error
		perms, err = strconv.ParseUint(v, 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid permissions value"))
			return
		}
	case float64:
		perms = uint64(v)
	default:
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid permissions claim type"))
		return
    }

	if claims["username"] == reqBody.Username || reqBody.Username == "admin" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid user to modify permissions for"))
		return
	}

    currentPerms, err := db.GetUserPermissions(reqBody.Username)

    if err != nil {
        w.Header().Set("Content-Type", "text/plain")
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("error retrieving requested user permissions"))
        return
    }

    parsedPerms, err := strconv.ParseUint(reqBody.Permissions, 10, 64)

    if err != nil {
        w.Header().Set("Content-Type", "text/plain")
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("error parsing permissions string"))
        return
    }

    // Debug prints to understand the issue
    fmt.Printf("User perms (who is modifying): %b (%d)\n", perms, perms)
    fmt.Printf("Target user current perms: %b (%d)\n", currentPerms, currentPerms)
    fmt.Printf("New perms being set: %b (%d)\n", parsedPerms, parsedPerms)

    unauthorizedPerms := parsedPerms &^ perms

    if unauthorizedPerms != 0 {
        fmt.Printf("Unauthorized perms: %b (%d)\n", unauthorizedPerms, unauthorizedPerms)
        w.Header().Set("Content-Type", "text/plain")
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("unauthorized permissions modified"))
        return
    }

	err = db.SetUserPerms(reqBody.Username, parsedPerms)

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error updating user permissions"))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user perms updated"))
}
