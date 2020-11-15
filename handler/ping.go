package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/920oj/Kisaragi/client"
	"github.com/bwmarrin/discordgo"
)

// PingHandler pingコマンド用Handler
func PingHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch {
	case strings.HasPrefix(m.Content, "!ping"):
		now := time.Now()
		timestamp, err := m.Timestamp.Parse()
		if err != nil {
			fmt.Println(err)
		}
		client.SendMessage(s, m.ChannelID, "Pong!: "+now.Sub(timestamp).String())
	}
}
