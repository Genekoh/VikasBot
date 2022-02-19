package commands

import (
	"github.com/bwmarrin/discordgo"
)

type PingPongCommand struct{}

func (p PingPongCommand) GetCommandName() string {
	return "PingPong"
}

func (p PingPongCommand) GetCommandCallers() []string {
	return []string{"ping", "pong"}
}

func (p PingPongCommand) GetCommandDesc() string {
	return "Replies back with \"ping\" or \"pong\"."
}

func (p PingPongCommand) Exec(sess *discordgo.Session, msg *discordgo.Message, _ []string, caller string) error {
	var responseString string
	if caller == "ping" {
		responseString = "pong"
	} else {
		responseString = "ping"
	}

	_, err := sess.ChannelMessageSend(msg.ChannelID, responseString)
	if err != nil {
		return err
	}

	return nil
}
