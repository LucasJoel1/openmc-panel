package globals

import (
	"encoding/json"
	"io"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/shirou/gopsutil/v4/process"
)

var discord *discordgo.Session
var serverRunning bool
var serverSettings *ServerSettingsStruct = nil
var serverPipes *PipesStruct = &PipesStruct{}
var lockServer bool
var startTime int64 = 0
var proc *process.Process = nil
var versions ServerVersions


func InitDiscord(token string) *discordgo.Session {
	var err error
	if discord != nil {
		return nil
	}
	discord, err = discordgo.New("Bot " + token)
	if err != nil {
		return nil
	}
	return discord
}

func SetServerVersions(minecraftVersion string, modloader string , modloaderVersion string) {
	versions.MinecraftVersion = minecraftVersion
	versions.Modloader = modloader
	versions.ModloaderVersion = modloaderVersion
}

func GetServerVersions() ServerVersions {
	return versions
}

func GetServerSettings() *ServerSettingsStruct {
	if serverSettings == nil {
		data, err := os.ReadFile("server.json")
		if err != nil {
			return serverSettings
		}

		var settings ServerSettingsStruct
		err = json.Unmarshal(data, &settings)
		if err != nil {
			return serverSettings
		}
		serverSettings = &settings
	}
	return serverSettings
}

func GetDiscord() *discordgo.Session { return discord }
func SetServerRunning(running bool)  { serverRunning = running }
func GetServerRunning() bool         { return serverRunning }
func SetPipes(_stdout io.ReadCloser, _stdin io.WriteCloser) {
	serverPipes.Stdout = _stdout;
	serverPipes.Stdin = _stdin
}
func GetPipes() PipesStruct {
	return *serverPipes
}
func LockServer(lock bool) { lockServer = lock }
func GetLocked() bool { return lockServer }
func SetStartTime(_startTime int64) { startTime = _startTime }
func GetStartTime() int64 { return startTime }
func SetProcess(_proc *process.Process) { proc = _proc }
func GetProcess() *process.Process { return proc }