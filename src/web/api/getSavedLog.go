package api

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func GetSavedLog(w http.ResponseWriter, r *http.Request) {
	logID := r.URL.Query().Get("logID")

	w.Header().Set("Content-Type", "application/gzip")

	dir, err := os.ReadDir("./logs")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
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
                w.WriteHeader(http.StatusInternalServerError)
            }
            return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
