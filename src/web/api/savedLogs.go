package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type logType struct {
    Name string `json:"name"`
    Size int64 `json:"size"`
}

func GetSavedLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dir, err := os.ReadDir("./logs")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "error reading logs directory",
			"errorCode": 1,
		})

		return
	}

	var logs []logType

	for _, log := range dir {
		logInfo, err := log.Info()
		if err != nil {
			fmt.Println(err)
			continue
		}
		item := logType{
			Name: strings.Split(log.Name(), ".")[0],
			Size: logInfo.Size(),
		}
		logs = append(logs, item)
		test, _ := log.Info()
		test.Size()
	}

	json.NewEncoder(w).Encode(map[string][]logType{
		"logs": logs,
	})
}
