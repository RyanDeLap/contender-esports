package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type ChannelMessageInfo struct {
	Msg *discordgo.MessageSend
	ChannelID string
}

func RunDiscordBot(c chan ChannelMessageInfo, apiKey string) {
	dg, err := discordgo.New("Bot " + apiKey)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return

	}
	handleIncomingMessages(c, dg)

	dg.Close()
}

func handleIncomingMessages(c chan ChannelMessageInfo, session *discordgo.Session) {
	for {
		data := <- c

		fmt.Println("Data", data.Msg.Content)

		if data.ChannelID == "" {
			break
		}
		_, err := session.ChannelMessageSendComplex(data.ChannelID, data.Msg)

		if err != nil {
			fmt.Print(err)
		}

	}
}

func messageCreate(session *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == session.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, "!setup") {
		session.ChannelMessageSend(m.ChannelID, "WIP")
	}
}

