package api

import (
	"net/http"
	"openmc-panel/src/globals"
)

func HandleServerLogsWS(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	connPtr := globals.AppendConnection(conn)

	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				conn.Close()
				globals.DeleteConnection(connPtr)
				return
			}
		}
	}()
}