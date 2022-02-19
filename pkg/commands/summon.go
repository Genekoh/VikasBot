package commands

import "github.com/bwmarrin/discordgo"

type SummonCommand struct{}

func (s SummonCommand) GetCommandName() string {
	return "Summon"
}

func (s SummonCommand) GetCommandCallers() []string {
	return []string{"summon"}
}

func (s SummonCommand) GetCommandDesc() string {
	return "Summons the user to the voice channel"
}

func (s SummonCommand) Exec(sess *discordgo.Session, msg *discordgo.Message, _ []string, _ string) error {
	_, err := joinVoiceChannel(sess, msg)
	if err != nil {
		return err
	}

	return nil
}

func joinVoiceChannel(sess *discordgo.Session, msg *discordgo.Message) (*discordgo.VoiceConnection, error) {
	isInChannel, vcId, err := isUserInVoiceChannel(sess, msg)
	if err != nil {
		return nil, err
	}
	if !isInChannel {
		_, err = sess.ChannelMessageSend(msg.ChannelID, "You have to be in a voice channel to call this command.")
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	vcConn, err := sess.ChannelVoiceJoin(msg.GuildID, vcId, false, false)
	if err != nil {
		return nil, err
	}

	return vcConn, nil
}

func isUserInVoiceChannel(sess *discordgo.Session, msg *discordgo.Message) (bool, string, error) {
	guild, err := sess.State.Guild(msg.GuildID)
	if err != nil {
		return false, "", err
	}

	isInChannel := false
	chanId := ""

	for _, vs := range guild.VoiceStates {
		if vs.UserID == msg.Author.ID {
			isInChannel = true
			chanId = vs.ChannelID
		}
	}

	return isInChannel, chanId, nil
}
