package sharedListeners

import (
	"fmt"
	"openmc-panel/src/globals"
	"strings"
)

func OnConnect(line string) {
	parts := strings.Split(line, ": ")
	usernameAndIP := strings.Split(parts[1], "[/")
	username := usernameAndIP[0]
	ip := strings.Split(usernameAndIP[1], ":")[0]
	_, err := globals.AddPlayer(username, ip)
	if err != nil {
		fmt.Println(err)
	}
}

func OnDisconnect(line string) {
	parts := strings.Split(line, ": ")
	username := strings.Split(parts[1], " ")[0]
	globals.DeletePlayer(username)
}