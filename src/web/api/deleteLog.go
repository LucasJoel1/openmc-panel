package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DeleteLog(w http.ResponseWriter, r *http.Request) {
	logID := r.URL.Query().Get("logID")

	if logID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if strings.ContainsAny(logID, "/\\..") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dir, err := os.ReadDir("./logs")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}

	for _, log := range dir {
		if strings.Split(log.Name(), ".")[0] == logID {
			filePath := filepath.Join("logs", log.Name())
			absLogs, _ := filepath.Abs("logs")
			absFile, _ := filepath.Abs(filePath)

            if !strings.HasPrefix(absFile, absLogs) {
                w.WriteHeader(http.StatusForbidden)
                return
            }

			err := os.Remove(filePath)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                return
            }
            w.WriteHeader(http.StatusOK)
            return
		}
	}
	
    w.WriteHeader(http.StatusNotFound)
}