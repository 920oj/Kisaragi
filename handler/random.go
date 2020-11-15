package handler

import (
	"github.com/920oj/Kisaragi/client"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strings"
	"time"
)

// RandomHandler Handler for Random command
func RandomHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch {
	case strings.HasPrefix(m.Content, "!random"):
		t := strings.Replace(m.Content, "!random ", "", 1)
		words := strings.Split(t, ",")
		rand.Seed(time.Now().UnixNano())
		if len(words) > 1 {
			client.SendMessage(s, m.ChannelID, words[rand.Intn(len(words))])
		} else {
			client.SendMessage(s, m.ChannelID, "You must specify at least two words.")
		}
	}
}
