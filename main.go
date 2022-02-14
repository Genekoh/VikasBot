package main

import (
	"fmt"
	"github.com/Genekoh/VikasBot/pkg/bot"
	"github.com/Genekoh/VikasBot/pkg/env"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var commandList = []bot.Command{
	bot.PingCommand{},
	bot.PongCommand{},
}

func main() {
	environ, err := env.GetEnviron("")
	if err != nil {
		log.Fatal(err)
	}
	AuthKey := environ["API_KEY"]
	if AuthKey == "" {
		log.Fatal("Bot Authentication Key not specified")
	}

	dg, err := discordgo.New("Bot " + AuthKey)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Println("error opening connection: ", err)
		return
	}
	defer func() {
		err = dg.Close()
		if err != nil {
			log.Fatal("unable to close: ", err)
		}
	}()

	blockUntilInterrupt()
}

func messageCreate(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	fmt.Println("message received")
	if sess.State.User.ID == msg.Author.ID {
		return
	}
	if !bot.StartWithPrefix(msg.Message) {
		return
	}

	handleCommand(sess, msg.Message)
}

func handleCommand(sess *discordgo.Session, msg *discordgo.Message) {
	a := bot.ParseArgs(msg)
	userCommandName := a[0]
	args := a[1:]

	for _, c := range commandList {
		name := c.GetCommandName()
		if name != userCommandName {
			continue
		}

		err := c.Exec(sess, msg, args)
		if err != nil {
			log.Println("err: ", err)
		}
	}
}

func blockUntilInterrupt() {
	fmt.Println("Bot Succesfully Ran")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
