package globals

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"time"
	"github.com/prometheus-community/pro-bing"
	"github.com/bwmarrin/discordgo"
	"github.com/disintegration/imaging"
)

var players = make(map[string]*Player)
var playerTicker *time.Ticker = time.NewTicker(10000 * time.Millisecond)

func StartPingRetrieval() {
	for range playerTicker.C {
		if len(players) > 0 {
			for _, player := range players {
				if player == nil {
					continue
				}

				pinger, err := probing.NewPinger(player.IP)
				if err != nil {
					fmt.Println("Ping error:", err)
					player.Ping = -1
					continue
				}

				pinger.SetPrivileged(true)
				pinger.Count = 3
				pinger.Timeout = 3 * time.Second

				err = pinger.Run()
				if err != nil {
					pinger, err = probing.NewPinger(player.IP)
					if err != nil {
						player.Ping = -1
						continue
					}

					pinger.SetPrivileged(false)
					pinger.Count = 3
					pinger.Timeout = 3 * time.Second

					err = pinger.Run()
					if err != nil {
						player.Ping = -1
						continue
					}
				}

				stats := pinger.Statistics()
				player.Ping = int(stats.AvgRtt.Milliseconds())
			}
		}
	}
}

func AddPlayer(username string, ip string) (*Player, error) {
	if username == "" {
		return nil, fmt.Errorf("username does not exist")
	}

	res, err := http.Get("https://api.mojang.com/users/profiles/minecraft/" + username)

	if err != nil {
		return nil, fmt.Errorf("error fetching user id")
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("error fetching user id")
	}

	res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("user ID request failed")
	}

	var idRequest IDRequest
	err = json.Unmarshal(body, &idRequest)

	if err != nil {
		return nil, fmt.Errorf("error decoding id request")
	}

	res, err = http.Get("https://sessionserver.mojang.com/session/minecraft/profile/" + idRequest.ID)

	if err != nil {
		return nil, fmt.Errorf("error fetching user skin")
	}

	body, err = io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("error parsing skin data")
	}

	res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("skin request failed")
	}

	var skinRequest SkinRequest
	err = json.Unmarshal(body, &skinRequest)

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error parsing skin request")
	}

	userSkinData, err := base64.StdEncoding.DecodeString(skinRequest.Properties[0].Value)

	if err != nil {
		return nil, fmt.Errorf("error decoding skin value")
	}

	var skinData SkinData

	err = json.Unmarshal(userSkinData, &skinData)

	if err != nil {
		return nil, fmt.Errorf("error parsing skin data")
	}

	res, err = http.Get(skinData.Textures.Skin.Url)

	if err != nil {
		return nil, fmt.Errorf("error retrieving skin")
	}

	body, err = io.ReadAll(res.Body)

	res.Body.Close()

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("cannot read skin")
	}

	img, _, err := image.Decode(bytes.NewReader(body))

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error processing image")
	}

	crop := image.Rect(8, 8, 16, 16)

	cropped := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(crop)

	resizedImg := imaging.Resize(cropped, 256, 256, imaging.NearestNeighbor)

	var buf bytes.Buffer

	err = png.Encode(&buf, resizedImg)

	if err != nil {
		return nil, err
	}

	b64Image := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	chatChannel := os.Getenv("CHAT_CHANNEL")

	webhook, err := GetDiscord().WebhookCreate(chatChannel, username, b64Image)

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("cannot create webhook")
	}

	player := Player{
		UUID: idRequest.ID,
		Username: username,
		Webhook: *webhook,
		IP: ip,
		Ping: 0,
		JoinedAt: int(time.Now().UnixMilli() / 1000),
	}

	players[username] = &player

	discord.WebhookExecute(webhook.ID, webhook.Token, true, &discordgo.WebhookParams{
		Content: username + " has joined the game",
	})

	return &player, nil
}

func GetPlayer(username string) *Player {
	return players[username]
}

func DeletePlayer(username string) error {
	webhook := players[username].Webhook

	discord.WebhookExecute(webhook.ID, webhook.Token, true, &discordgo.WebhookParams{
		Content: username + " has left the game",
	})

	err := GetDiscord().WebhookDelete(webhook.ID)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("could not delete webhook")
	}
	
	players[username] = nil

	return nil
}

func GetPlayers() []*Player {
	var playersArr []*Player

	for _, player := range players {
		if player != nil {
			playersArr = append(playersArr, player)
		}
	}

	return playersArr
}