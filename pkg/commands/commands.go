package commands

import (
	"github.com/bwmarrin/discordgo"
)

var commandList = []Command{
	PingPongCommand{},
	PrefixCommand{},
	SummonCommand{},
	PlayCommand{},
	HelpCommand{},
}

type Command interface {
	GetCommandName() string
	GetCommandCallers() []string
	GetCommandDesc() string
	Exec(sess *discordgo.Session, msg *discordgo.Message, args []string, caller string) error
}

func MatchCommandCallers(match string) Command {
	var command Command
	for _, c := range commandList {
		for _, caller := range c.GetCommandCallers() {
			if caller != match {
				continue
			}
			command = c
		}

	}

	return command
}

func MatchCommandName(match string) Command {
	var command Command
	for _, c := range commandList {
		if match == c.GetCommandName() {
			command = c
		}
	}

	return command
}
