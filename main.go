package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	// TokenPrefix Tokenは先頭にBotという文字が必要
	TokenPrefix = "Bot "
	// BotName = "<@777372032333119509>"
)

func main() {
	// インスタンス生成
	discord, err := discordgo.New(TokenPrefix + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	discord.AddHandler(PingHandler)

	// WebSocket開始
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	// サーバー開始時の処理
	fmt.Println("Start Listening...\nPress CTRL-C to exit.")

	// bot終了用の処理
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discord.Close()
}
