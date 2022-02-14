package bot

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

const defaultPrefix = "!!"

func StartWithPrefix(msg *discordgo.Message) bool {
	m := msg.Content
	prefix := defaultPrefix

	if m[0:len(prefix)] == prefix {
		return true
	}

	return false
}

func ParseArgs(msg *discordgo.Message) []string {
	s := removePrefix(msg.Content)
	args := strings.Split(s, " ")

	trimmedArgs := make([]string, len(args))
	for i, a := range args {
		trimmedArgs[i] = strings.TrimSpace(a)
	}

	return trimmedArgs
}

func removePrefix(m string) string {
	prefix := defaultPrefix

	return m[len(prefix):]
}
