package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type HelpCommand struct{}

func (h HelpCommand) GetCommandName() string {
	return "Help"
}

func (h HelpCommand) GetCommandCallers() []string {
	return []string{"help", "h"}
}

func (h HelpCommand) GetCommandDesc() string {
	return "Replies back with a description of the specified command."
}

func (h HelpCommand) Exec(sess *discordgo.Session, msg *discordgo.Message, args []string, caller string) error {
	var response string
	if len(args) == 0 || args[0] == "all" {
		for i, command := range commandList {
			if i == 0 {
				response += formatCommandDesc(command)
			} else {
				d := formatCommandDesc(command)
				response += fmt.Sprintf("\n\n%s", d)
			}
		}

		_, err := sess.ChannelMessageSend(msg.ChannelID, response)
		if err != nil {
			return err
		}

		return nil
	}

	if command := MatchCommandName(args[0]); command != nil {
		response = formatCommandDesc(command)
	} else if command := MatchCommandCallers(args[0]); command != nil {
		response = formatCommandDesc(command)
	} else {
		response = "Give a valid command."
	}

	_, err := sess.ChannelMessageSend(msg.ChannelID, response)
	if err != nil {
		return err
	}

	return nil
}

func formatCommandDesc(command Command) string {
	var callers string
	for i, caller := range command.GetCommandCallers() {
		if i == 0 {
			callers += caller
		} else {
			callers += fmt.Sprintf(", %s", caller)
		}
	}

	return fmt.Sprintf(
		"Name: %s\nCallers: %s\nDescription: %s", command.GetCommandName(), callers, command.GetCommandDesc(),
	)
}
