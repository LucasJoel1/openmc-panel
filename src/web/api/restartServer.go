package api

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"openmc-panel/src/globals"
	sharedListeners "openmc-panel/src/listeners/share"
	"openmc-panel/src/server"
)

func HandleRestartReq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !globals.GetServerRunning() {
		w.WriteHeader(http.StatusConflict)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "server not running",
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

	writer := bufio.NewWriter(globals.GetPipes().Stdin)
	writer.WriteString("stop" + "\n")
	writer.Flush()

	globals.SetServerRunning(false)

	_, err := server.StartServer()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   err.Error(),
			"errorCode": 3,
		})
		return
	}

	ctx := context.Background()

	go func() {
		sharedListeners.LogsListener(ctx)
	}()

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Server Restarted Successfully",
		"errorCode": 0,
	})
}