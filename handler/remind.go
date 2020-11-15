package handler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/920oj/Kisaragi/client"
	"github.com/bwmarrin/discordgo"
)

// RemindHandler Handler for Remind command
func RemindHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch {
	case strings.HasPrefix(m.Content, "!remind"):
		t := strings.Replace(m.Content, "!remind ", "", 1)
		msg, err := createRemindMessage(t)
		if err != nil {
			fmt.Println("hoge")
		}
		client.SendMessage(s, m.ChannelID, msg)
	}
}

func createRemindMessage(t string) (string, error) {
	timeReg, _ := regexp.Compile(`(^.*?)後に.*`)
	contentReg, _ := regexp.Compile(`^.*?後に(.*)を通知$`)

	time := timeReg.ReplaceAllString(t, "$1")
	content := contentReg.ReplaceAllString(t, "$1")

}

func parseTime(t string) (string, error) {
	reg, _ := regexp.Compile(``)
}
