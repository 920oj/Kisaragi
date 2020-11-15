package client

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func SendMessage(s *discordgo.Session, ChannelID string, msg string) {
	_, err := s.ChannelMessageSend(ChannelID, msg)

	log.Println(">>>" + msg)
	if err != nil {
		log.Println("Error: ", err)
	}
}
