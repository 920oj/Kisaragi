package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"github.com/920oj/Kisaragi/handler"
)

var (
	// TokenPrefix Token must be prefixed with the character "Bot"
	TokenPrefix = "Bot "
	stopChannel chan bool
)

func main() {
	loadEnv()

	discord, err := discordgo.New()
	discord.Token = TokenPrefix + os.Getenv("DISCORD_TOKEN")
	if err != nil {
		fmt.Println(err)
		return
	}

	stopChannel = make(chan bool)

	discord.AddHandler(handler.PingHandler)
	discord.AddHandler(handler.RandomHandler)
	discord.AddHandler(handler.RemindHandler)

	// Open Websocket
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Start Listening...\nPress CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	discord.Close()
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
