package globals

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"time"

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

func SetServerVersions(minecraftVersion string, modloader string, modloaderVersion string) {
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

func SetServerRunning(running bool) {
	if !running {
		lockServer = true
		go func() {
			if proc == nil {
				lockServer = false
				serverRunning = false
				return
			}
			for {
				isRunning, _ := proc.IsRunning()
				if !isRunning {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
			lockServer = false
			serverRunning = false

			if serverSettings == nil || serverSettings.Path == "" {
				log.Println("Server settings not initialized or path is empty")
				return
			}

			dataTime := time.Now().UnixMilli()
			srcPath := path.Join(serverSettings.Path, "logs", "latest.log")
			dstPath := path.Join("./logs", strconv.FormatInt(dataTime, 10)+".log")

			srcFile, err := os.Open(srcPath)
			if err != nil {
				log.Println(err)
				return
			}
			defer srcFile.Close()

			dstFile, err := os.Create(dstPath + ".gz")
			if err != nil {
				log.Println(err)
				return
			}
			defer dstFile.Close()

			gzipWriter := gzip.NewWriter(dstFile)
			defer gzipWriter.Close()

			_, err = io.Copy(gzipWriter, srcFile)
			if err != nil {
				fmt.Println(err)
			}
		}()
	} else {
		serverRunning = running
	}
}

func GetServerRunning() bool { return serverRunning }
func SetPipes(_stdout io.ReadCloser, _stdin io.WriteCloser) {
	serverPipes.Stdout = _stdout
	serverPipes.Stdin = _stdin
}
func GetPipes() PipesStruct {
	return *serverPipes
}
func LockServer(lock bool)              { lockServer = lock }
func GetLocked() bool                   { return lockServer }
func SetStartTime(_startTime int64)     { startTime = _startTime }
func GetStartTime() int64               { return startTime }
func SetProcess(_proc *process.Process) { proc = _proc }
func GetProcess() *process.Process      { return proc }
