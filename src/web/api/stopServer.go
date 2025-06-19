package api

import (
	"bufio"
	"encoding/json"
	"net/http"
	"openmc-panel/src/globals"
)


func HandleStopReq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !globals.GetServerRunning() {
		w.WriteHeader(http.StatusConflict)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "server is not running",
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

	globals.LockServer(true)
	defer globals.LockServer(false)

	globals.SetServerRunning(false)

	writer := bufio.NewWriter(globals.GetPipes().Stdin)
	writer.WriteString("stop" + "\n")
	writer.Flush()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Server Stopped Successfully",
		"errorCode": 0,
	})
}