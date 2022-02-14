package bot

import (
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	GetCommandName() string
	GetCommandDesc() string
	Exec(sess *discordgo.Session, msg *discordgo.Message, args []string) error
}

type PingCommand struct{}

func (ping PingCommand) GetCommandName() string {
	return "ping"
}

func (ping PingCommand) GetCommandDesc() string {
	return "Replies back with \"pong\""
}

func (ping PingCommand) Exec(sess *discordgo.Session, msg *discordgo.Message, _ []string) error {

	_, err := sess.ChannelMessageSend(msg.ChannelID, "pong")
	if err != nil {
		return err
	}

	return nil
}

type PongCommand struct{} // TODO: Refactor code so pong command is not a copy past of the ping command

func (pong PongCommand) GetCommandName() string {
	return "pong"
}

func (pong PongCommand) GetCommandDesc() string {
	return "Replies back with \"ping\""
}

func (pong PongCommand) Exec(sess *discordgo.Session, msg *discordgo.Message, _ []string) error {

	_, err := sess.ChannelMessageSend(msg.ChannelID, "ping")
	if err != nil {
		return err
	}

	return nil
}
