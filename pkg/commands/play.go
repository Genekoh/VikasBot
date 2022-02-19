package commands

import (
	"fmt"
	youtube "github.com/Genekoh/VikasBot/pkg/api"
	"github.com/Genekoh/VikasBot/pkg/audio"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

type PlayCommand struct{}

func (p PlayCommand) GetCommandName() string {
	return "Play"
}

func (p PlayCommand) GetCommandCallers() []string {
	return []string{"play", "p"}
}

func (p PlayCommand) GetCommandDesc() string {
	return "Joins the user's current voice channel and plays audio. Has an alias as p"
}

func (p PlayCommand) Exec(sess *discordgo.Session, msg *discordgo.Message, args []string, _ string) error {
	if len(args) == 0 {
		_, err := sess.ChannelMessageSend(msg.ChannelID, "You need to enter a name of a video.")
		if err != nil {
			return err
		}

		return nil
	}

	go func() {
		vc, err := joinVoiceChannel(sess, msg)
		if err != nil {
			log.Println(err)
			return
		}

		queryString := strings.Join(args, " ")
		searchResults, err := youtube.QueryVideos(queryString, 1)
		result := searchResults[0]
		s := fmt.Sprintf("%s\n%s\n%s", result.Snippet.Title, result.Snippet.ChannelTitle, result.Snippet.Description)
		_, err = sess.ChannelMessageSend(msg.ChannelID, s)
		if err != nil {
			log.Println(err)
			return
		}

		song := audio.NewSong(result.Snippet.Title, result.Snippet.ChannelTitle, result.Id.VideoId)

		vi := audio.GetVoiceInstance(msg.GuildID, vc)
		vi.QueueAdd(song)

	}()

	return nil
}
