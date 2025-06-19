package sharedListeners

import (
	"openmc-panel/src/globals"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func OnGameMessage(line string) {
	parts := strings.Split(strings.Split(line, ": <")[1],  "> ")
	username := parts[0]
	message := parts[1]

	webhook := globals.GetPlayer(username).Webhook

	globals.GetDiscord().WebhookExecute(webhook.ID, webhook.Token, true, &discordgo.WebhookParams{
		Content: message,
	})
}