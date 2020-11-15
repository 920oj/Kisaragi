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
	// TokenPrefix Tokenは先頭にBotという文字が必要
	TokenPrefix = "Bot "
	// BotName = "<@777372032333119509>"
	stopChannel chan bool
)

func main() {
	loadEnv()
	// インスタンス生成
	discord, err := discordgo.New()
	discord.Token = TokenPrefix + os.Getenv("DISCORD_TOKEN")
	if err != nil {
		fmt.Println(err)
	}

	stopChannel = make(chan bool)
	discord.AddHandler(handler.PingHandler)

	// WebSocket開始
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	// サーバー開始時の処理
	fmt.Println("Start Listening...\nPress CTRL-C to exit.")

	// bot終了用の処理
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	discord.Close()
}

// loadEnv 環境変数の読み込み
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
