package main

import (
	"fmt"
	"log"
	"openmc-panel/src/globals"
	"openmc-panel/src/listeners/bot"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"openmc-panel/src/web"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("COULD NOT LOAD .env FILE")
	}
	token := os.Getenv("DISCORD_TOKEN")

	if token == "" {
		log.Fatal("DISCORD_TOKEN is empty in .env")
	}

	app := globals.InitDiscord(token)

	if app == nil {
		log.Fatal("FAILED TO INITIALIZE DISCORD BOT")
	}

	app.AddHandler(botListners.OnMessage)

	app.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	err = app.Open()

	if err != nil {
		log.Fatal("Error opening Discord client ", err)
	}

	go web.StartWeb()
	go globals.StartPingRetrieval()
	go globals.StartPropertiesParser()

	globals.SetServerVersions(globals.GetServerType())

	fmt.Println("App online. Press CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	app.Close()
}