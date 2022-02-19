package bot

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

const DefaultPrefix = "!!"

func ParseArgs(msg *discordgo.Message, prefix string) []string {
	s := removePrefix(msg.Content, prefix)
	args := strings.Split(s, " ")

	trimmedArgs := make([]string, len(args))
	for i, a := range args {
		trimmedArgs[i] = strings.TrimSpace(a)
	}

	return trimmedArgs
}

func removePrefix(m string, prefix string) string {

	return m[len(prefix):]
}
