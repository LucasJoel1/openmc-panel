package globals

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var properties = make(map[string]Property)

// var bannedPlayers = make(map[string]string)
// var whitelistedPlayers = make(map[string]string)
var propertiesTicker *time.Ticker = time.NewTicker(5000 * time.Millisecond)

var commentsMatcher, _ = regexp.Compile(`^ *#.*\n$`)

var mcFabQuiltNeoFolderVersion, _ = regexp.Compile(`^([0-9]*\.?){3}$`)
var forgeFolderVersion, _ = regexp.Compile(`^(([0-9]*\.?){3}-?){2}$`)

func StartPropertiesParser() {
	for range propertiesTicker.C {

		for k := range properties {
			delete(properties, k)
		}

		serverPath := GetServerSettings().Path
		f, err := os.Open(path.Join(serverPath, "server.properties"))

		if err != nil {
			return
		}

		r := bufio.NewReader(f)

		for {
			line, err := r.ReadString('\n')
			if err != nil {
				break
			}

			if commentsMatcher.MatchString(line) {
				continue
			}

			line = strings.ReplaceAll(strings.ReplaceAll(line, " ", ""), "\r\n", "")
			parts := strings.Split(line, "=")

			key := parts[0]

			value := parts[1]

			prop := properties[key]
			if checkStringNumber(value) {
				if intValue, err := strconv.Atoi(value); err == nil {
					prop.Value = intValue
					properties[key] = prop
				} else {
					prop.Value = value
					properties[key] = prop
				}
			} else {
				switch value {
				case "true":
					prop.Value = true
				case "false":
					prop.Value = false
				default:
					prop.Value = value
				}
			}

			prop.lastUpdated = int(time.Now().UnixMilli() * 1000)
			properties[key] = prop
		}

		f.Close()
	}
}

func GetProperty(key string) Property {
	return properties[key]
}

func GetPropertyAsString(key string) (string, int) {
	prop := properties[key]
	if prop.Value == nil {
		return "", 0
	}
	return fmt.Sprintf("%v", prop.Value), prop.lastUpdated
}

func SetProperty(key string, value any) bool {
	prop := properties[key]
	if prop.Value == nil {
		return false
	}
	prop.Value = value
	prop.lastUpdated = int(time.Now().UnixMilli() * 1000)
	properties[key] = prop
	return true
}

func checkStringNumber(str string) bool {
	for _, val := range str {
		if val < '0' || val > '9' {
			return false
		}
	}
	return true
}

func GetServerType() (minecraftServerVersion string, modLoader string, modLoaderVersion string) {
	serverFolder := GetServerSettings().Path

	if dirs, err := os.ReadDir(path.Join(serverFolder, "libraries", "net", "fabricmc", "fabric-loader")); err == nil {
		modLoader = "Fabric"
		for _, v  := range dirs {
			if v.IsDir() && mcFabQuiltNeoFolderVersion.MatchString(v.Name()) {
				modLoaderVersion = v.Name()
				break;
			}
		}
		minecraftServerVersion = fabricBasedGetServerVersion(serverFolder)
	} else if dirs, err := os.ReadDir(path.Join(serverFolder, "libraries", "org", "quiltmc", "quilt-loader")); err == nil {
		modLoader = "Quilt"
		for _, v  := range dirs {
			if v.IsDir() && mcFabQuiltNeoFolderVersion.MatchString(v.Name()) {
				modLoaderVersion = v.Name()
				break;
			}
		}
		minecraftServerVersion = fabricBasedGetServerVersion(serverFolder)	
	} else if dirs, err := os.ReadDir(path.Join(serverFolder, "libraries", "net", "minecraftforge", "forge")); err == nil {
		modLoader = "Forge"
		for _, v := range dirs {
			if v.IsDir() && forgeFolderVersion.MatchString(v.Name()) {
				parts := strings.Split(v.Name(), "-")
				minecraftServerVersion = parts[0]
				modLoaderVersion = parts[1]
				break;
			}
		}
	} else if dirs, err := os.ReadDir(path.Join(serverFolder, "libraries", "net", "neoforged", "neoforge")); err == nil {
		modLoader = "NeoForge"
		for _, v := range dirs {
			if v.IsDir() && mcFabQuiltNeoFolderVersion.MatchString(v.Name()) {
				modLoaderVersion = v.Name()
				parts := strings.Split(v.Name(), ".")
				minecraftServerVersion = "1." + parts[0] + "." + parts[1]
				break;
			}
		}
	} else if _, err := os.ReadFile(path.Join(serverFolder, "paper.yml")); err == nil {
		modLoader = "paper"
		dirs, err := os.ReadDir(path.Join(serverFolder, "versions"))

		if err != nil {
			return
		}

		if len(dirs) == 1 {
			versionsDir, err := os.ReadDir(path.Join(serverFolder, "versions", dirs[0].Name()))
			if err != nil {
				return
			}
			if len(versionsDir) == 1 {
				r, err := zip.OpenReader(path.Join(serverFolder, "versions", dirs[0].Name(), versionsDir[0].Name()))

				if err != nil {
					return
				}

				defer r.Close()
				for _, f := range r.File {
					if f.Name == "META-INF/MANIFEST.MF" {
						rc, err := f.Open()
						if err != nil {
							continue
						}
						
						defer rc.Close()

						content, err := io.ReadAll(rc)

						if err != nil {
							continue
						}

						for _, line := range strings.Split(string(content), "\n") {
							if strings.HasPrefix(line, "Implementation-Version: ") {
								fullVersion := strings.Split(line, "Implementation-Version: ")[1]
								parts := strings.Split(fullVersion, "-")
								modLoaderVersion = parts[1]
								minecraftServerVersion = parts[0]
							}
						}
					}
				}
			}
		}
	}
	return
}

func fabricBasedGetServerVersion(serverFolder string) (version string) {
	dir, err := os.ReadDir(path.Join(serverFolder, "versions"))

	if err != nil {
		return
	}

	for _, v := range dir {
		if v.IsDir() && mcFabQuiltNeoFolderVersion.MatchString(v.Name()) {
			return v.Name()
		}
	}

	return
}