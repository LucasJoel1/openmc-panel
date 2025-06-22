package api

import (
	"fmt"
	"net/http"
	"openmc-panel/src/globals"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type player struct {
	Ping int    `json:"ping"`
	Time int    `json:"time"`
	Name string `json:"name"`
}

type serverInfoMessage struct {
	Status struct {
		State    int     `json:"state"`
		Uptime   int64   `json:"uptime"`
		TPS      float32 `json:"TPS"`
		CPU      float32 `json:"CPU"`
		RAMUsage float32 `json:"RAMUsage"`
		RAMAlloc float32 `json:"RamAlloc"`
	}

	Players struct {
		Online int      `json:"online"`
		Max    int      `json:"max"`
		List   []player `json:"list"`
	}

	Properties struct {
		Version          string `json:"version"`
		GameMode         string `json:"GameMode"`
		Difficulty       string `json:"Difficulty"`
		IP               string `json:"IP"`
		Port             int    `json:"Port"`
		ModLoader        string `json:"ModLoader"`
		ModLoaderVersion string `json:"ModLoaderVersion"`
	}
}

var ws = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleServerInfoWS(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	data := serverInfoMessage{}

	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				ticker.Stop()
				conn.Close()
				return
			}
		}
	}()

	for range ticker.C {
		err = setData(&data)
		if err != nil {
			break
		}

		if err := conn.WriteJSON(data); err != nil {
			break
		}
	}
}

func setData(data *serverInfoMessage) error {
	// SERVER PROPS
	if globals.GetServerRunning() {
		data.Status.State = 0
	} else if globals.GetLocked() && !globals.GetServerRunning() {
		data.Status.State = 1
	} else {
		data.Status.State = 2
	}

	if !globals.GetServerRunning() {
		data.Status.Uptime = 0
	} else {
		data.Status.Uptime = time.Now().Unix() - globals.GetStartTime()
	}

	proc := globals.GetProcess()

	if proc != nil {

		children, err := proc.Children()

		if err != nil {
			return fmt.Errorf("cannot retrieve children")
		}

		memory := 0

		for _, child := range children {
			meminfo, err := child.MemoryInfo()

			if err != nil {
				return fmt.Errorf("cannot get memory for a child process")
			}
			memory += int(meminfo.RSS)

		}

		data.Status.RAMUsage = float32(memory) / 1_024_000_000

		data.Status.TPS = 20.0 // todo

		cpuPercent, _ := proc.CPUPercent()
		data.Status.CPU = float32(cpuPercent)

		ramAlloc, _ := strconv.ParseFloat(strings.Split(globals.GetServerSettings().Memory[1], "G")[0], 32)
		data.Status.RAMAlloc = float32(ramAlloc)
	} else {
		data.Status.RAMUsage = 0.0
		data.Status.TPS = 0.0
		data.Status.CPU = 0.0
	}

	// PLAYER PROPS
	data.Players.Online = len(globals.GetPlayers())
	maxPlayers, _ := globals.GetPropertyAsString("max-players")
	data.Players.Max, _ = strconv.Atoi(maxPlayers)
	data.Players.List = data.Players.List[:0]
	for _, v := range globals.GetPlayers() {
		if v != nil {
			data.Players.List = append(data.Players.List, player{
				Name: v.Username,
				Ping: v.Ping,
				Time: int(time.Now().UnixMilli()/1000) - v.JoinedAt,
			})
		}
	}

	serverVersions := globals.GetServerVersions()

	// PROPERTIES
	data.Properties.Version = serverVersions.MinecraftVersion
	data.Properties.GameMode, _ = globals.GetPropertyAsString("gamemode")
	data.Properties.Difficulty, _ = globals.GetPropertyAsString("difficulty")
	data.Properties.IP = "192.168.0.1" // todo
	port, _ := globals.GetPropertyAsString("query.port")
	data.Properties.Port, _ = strconv.Atoi(port)
	data.Properties.ModLoader = serverVersions.Modloader
	data.Properties.ModLoaderVersion = serverVersions.ModloaderVersion
	return nil
}
