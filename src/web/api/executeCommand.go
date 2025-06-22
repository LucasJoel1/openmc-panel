package api

import (
	"bufio"
	"encoding/json"
	"net/http"
	"openmc-panel/src/globals"
)

func ExecuteCommand(w http.ResponseWriter, r *http.Request) {
	command := r.URL.Query().Get("command")

	writer := bufio.NewWriter(globals.GetPipes().Stdin)

	_, err := writer.WriteString(command + "\n")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   err.Error(),
			"errorCode": 1,
		})
		return
	}

	err = writer.Flush()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   err.Error(),
			"errorCode": 1,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Command sent successfully",
		"errorCode": 0,
	})
}