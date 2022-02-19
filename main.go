package main

import (
	"fmt"
	"github.com/Genekoh/VikasBot/pkg/bot"
	"github.com/Genekoh/VikasBot/pkg/commands"
	"github.com/Genekoh/VikasBot/pkg/database"
	"github.com/Genekoh/VikasBot/pkg/env"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	defer handleError()
	environ, err := env.GetEnviron("")
	if err != nil {
		log.Panic("error getting environment variables: ", err)
	}
	AuthKey := environ["DISCORD_AUTH_KEY"]
	if AuthKey == "" {
		log.Panic("Bot Authentication Key not specified")
	}

	dg, err := discordgo.New("Bot " + AuthKey)
	if err != nil {
		log.Panic("error creating Discord session: ", err)
	}

	err = database.ConnectDatabase()
	if err != nil {
		log.Panic("error connecting to database", err)
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Panic("error opening connection: ", err)
	}

	defer func() {
		err = dg.Close()
		if err != nil {
			log.Panic("unable to close: ", err)
		}
	}()

	blockUntilInterrupt()
}

func messageCreate(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	config, err := database.GetServerConfig(msg.GuildID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if sess.State.User.ID == msg.Author.ID {
		return
	}
	if !strings.HasPrefix(msg.Message.Content, config.Prefix) {
		return
	}
	fmt.Println("received command")

	handleCommand(sess, msg.Message, config)
}

func handleCommand(sess *discordgo.Session, msg *discordgo.Message, conf database.ServerConfig) {
	a := bot.ParseArgs(msg, conf.Prefix)
	userCall := a[0]
	args := a[1:]

	c := commands.MatchCommandCallers(userCall)
	if c == nil {
		return
	}

	err := c.Exec(sess, msg, args, userCall)
	if err != nil {
		log.Println("err: ", err)
	}
}

func blockUntilInterrupt() {
	fmt.Println("Bot Successfully Ran")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func handleError() {
	err := recover()
	log.Fatal(err)
}
