package web

import (
	"fmt"
	"log"
	"net/http"
	"openmc-panel/src/web/api"
)

func StartWeb() {
	fs := http.FileServer(http.Dir("./src/web/frontend"))

	http.Handle("/", fs)
	http.HandleFunc("/api/startServer", api.HandleStartReq)
	http.HandleFunc("/api/stopServer", api.HandleStopReq)
	http.HandleFunc("/api/restartServer", api.HandleRestartReq)
	http.HandleFunc("/api/ws/serverInfo", api.HandleServerInfoWS)
	http.HandleFunc("/api/ws/serverLogs", api.HandleServerLogsWS)

	fmt.Println("Serving on http://127.0.0.1")
	err := http.ListenAndServe(":80", nil)

	if err != nil {
		log.Fatal(err)
	}
}