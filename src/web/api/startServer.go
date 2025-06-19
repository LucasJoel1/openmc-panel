package api

import (
	"context"
	"encoding/json"
	"net/http"
	"openmc-panel/src/globals"
	"openmc-panel/src/server"
	"openmc-panel/src/listeners/share"
)

func HandleStartReq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	if globals.GetServerRunning() {
		w.WriteHeader(http.StatusConflict)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "server already running",
			"errorCode": 1,
		})
		return
	}

	if globals.GetLocked() {
		w.WriteHeader(http.StatusConflict)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "server is in locked state",
			"errorCode": 2,
		})
		return
	}

	_, err := server.StartServer()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   err.Error(),
			"errorCode": 2,
		})
		return
	}

	ctx := context.Background()

	go func() {
		sharedListeners.LogsListener(ctx)
	}()

	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Server Started Successfully",
		"errorCode": 0,
	})
}
