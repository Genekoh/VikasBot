package commands

import (
	"fmt"
	"github.com/Genekoh/VikasBot/pkg/database"
	"github.com/bwmarrin/discordgo"
)

type PrefixCommand struct{}

func (p PrefixCommand) GetCommandName() string {
	return "Prefix"
}

func (p PrefixCommand) GetCommandCallers() []string {
	return []string{"prefix"}
}

func (p PrefixCommand) GetCommandDesc() string {
	return "Change prefix the commands will be listening for in this server."
}

func (p PrefixCommand) Exec(sess *discordgo.Session, msg *discordgo.Message, args []string, caller string) error {
	if len(args) == 0 {
		_, err := sess.ChannelMessageSend(msg.ChannelID, "No Prefix Specified")
		if err != nil {
			return err
		}
		return nil
	}

	serverConfig, err := database.GetServerConfig(msg.GuildID)
	if err != nil {
		return err
	}

	serverConfig.Prefix = args[0]
	database.UpdateServerConfig(serverConfig)

	res := fmt.Sprintf("Prefix now change to %s", args[0])
	_, err = sess.ChannelMessageSend(msg.ChannelID, res)
	if err != nil {
		return err
	}
	
	return nil
}
