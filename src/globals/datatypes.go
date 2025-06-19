package globals

import (
	"io"

	"github.com/bwmarrin/discordgo"
)

type ServerSettingsStruct struct {
	ServerName string   `json:"serverName"`
	JarName    string   `json:"jarName"`
	Path       string   `json:"path"`
	Memory     []string `json:"memory"`
	GUI        bool     `json:"gui"`
}

type PipesStruct struct {
	Stdout io.ReadCloser
	Stdin  io.WriteCloser
}

type Player struct {
	Username string
	UUID     string
	Webhook  discordgo.Webhook
	IP       string
	Ping     int
	JoinedAt int
}

type IDRequest struct {
	ID string `json:"id"`
}

type SkinRequest struct {
	Properties []SkinPropertiesRequest `json:"properties"`
}

type SkinPropertiesRequest struct {
	Value string `json:"value"`
}

type SkinData struct {
	Textures SkinTextures `json:"textures"`
}

type SkinTextures struct {
	Skin SkinTexturesSKIN `json:"SKIN"`
}

type SkinTexturesSKIN struct {
	Url string `json:"url"`
}

type Property struct {
	Value any
	lastUpdated int
}

type Update struct {
	
}

type ServerVersions struct {
	MinecraftVersion string
	Modloader string
	ModloaderVersion string
}