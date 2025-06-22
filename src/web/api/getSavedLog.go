package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetSavedLog(w http.ResponseWriter, r *http.Request) {
	logID := r.URL.Query().Get("logID")

	if logID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if strings.ContainsAny(logID, "/\\..") {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/gzip")

	dir, err := os.ReadDir("./logs")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, log := range dir {
		if strings.Split(log.Name(), ".")[0] == logID {
			filePath := "./logs/" + log.Name()

			file, err := os.Open(filePath)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer file.Close()

			_, err = io.Copy(w, file)

            if err != nil {
                fmt.Println("Error streaming log file")
            }
            return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
