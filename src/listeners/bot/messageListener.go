package botListners

import (
	"bufio"
	"fmt"
	"openmc-panel/src/globals"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func OnMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	bot := globals.GetDiscord()

	if bot == nil {
		return
	}

	if message == nil {
		return
	}

	if message.Author.ID == bot.State.User.ID {
		return
	}

	if strings.HasPrefix(message.Content, "!") {
		command := strings.Split(message.Content, "!")[1]
		_, err := session.ChannelMessageSendReply(message.ChannelID, CommandHandler(command), (*message).Reference())
		if err != nil {
			fmt.Println(err)
		}
	} else if message.ChannelID == os.Getenv("CONSOLE_CHANNEL") && globals.GetServerRunning() && message.Author.ID != bot.State.User.ID {
		writer := bufio.NewWriter(globals.GetPipes().Stdin)
		writer.WriteString(message.Content + "\n")
		writer.Flush()
	} else if message.ChannelID == os.Getenv("CHAT_CHANNEL") && globals.GetServerRunning() && message.WebhookID == "" {
		writer := bufio.NewWriter(globals.GetPipes().Stdin)
		writer.WriteString(`tellraw @a [{"text":"[` + message.Author.Username + `] ` + `" ,"color":"blue"},{"text": "` + message.Content + `","color":"white"}]` + "\n")
		writer.Flush()
	}
}