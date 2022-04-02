package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// PingHandler pingコマンド用Handler
func PingHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ボットユーザからのメッセージは何もしない
	if m.Author.ID == s.State.User.ID {
		return
	}

	// "ping" というメッセージが来たら、反応速度を測って "pong" とともに返す
	if m.Content == "ping" {
		now := time.Now()
		timestamp := m.Timestamp
		p := timestamp.Sub(now).String()
		s.ChannelMessageSend(m.ChannelID, "pong: "+p)
	}
}
