package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	twitterURLPrefix = "https://twitter.com"
	twitterAPIURL    = "https://api.twitter.com/2/tweets"
	twitterQuery     = "&tweet.fields=id,text,author_id&media.fields=media_key,duration_ms,height,preview_image_url,type,url,width,public_metrics,alt_text&expansions=attachments.media_keys"
)

type twitterAPIRes struct {
	Data []struct {
		ID          string `json:"id"`
		AuthorID    string `json:"author_id"`
		Text        string `json:"text"`
		Attachments struct {
			MediaKeys []string `json:"media_keys"`
		} `json:"attachments"`
	} `json:"data"`
	Includes struct {
		Media []struct {
			Width    int    `json:"width"`
			URL      string `json:"url"`
			Height   int    `json:"height"`
			MediaKey string `json:"media_key"`
			Type     string `json:"type"`
		} `json:"media"`
	} `json:"includes"`
}

func downloadTwitterImgHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ボットユーザからのメッセージは何もしない
	if m.Author.ID == s.State.User.ID {
		return
	}

	// TwitterのURLから始まらない場合は何もしない
	if !strings.HasPrefix(m.Content, twitterURLPrefix) {
		return
	}

	tweetURL := m.Content
	headers := http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", os.Getenv("TW_BEARER_TOKEN"))},
	}

	// ツイートのIDを取得する
	rex := regexp.MustCompile(`([^\/]+$)`)
	tweetID := rex.FindString(tweetURL)

	// Twitter APIのリクエストURLを作成する
	ajaxURL := fmt.Sprintf("%s?ids=%s%s", twitterAPIURL, tweetID, twitterQuery)

	// Twitter APIを叩く
	fmt.Println(ajaxURL)
	twitterAPIBytes, err := RequestAjax(ajaxURL, headers)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Error: Cannot Download Twitter API Response.")
		return
	}

	var twitterAPIData twitterAPIRes
	json.Unmarshal(twitterAPIBytes, &twitterAPIData)

	fmt.Println(twitterAPIData)
	if len(twitterAPIData.Includes.Media) < 1 {
		// 画像が紐付いていないとき
		s.ChannelMessageSend(m.ChannelID, "Error: Images not found.")
		return
	} else if len(twitterAPIData.Includes.Media) == 1 {
		// 画像が1枚のとき
		imgURL := twitterAPIData.Includes.Media[0].URL
		u, err := url.Parse(imgURL)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error: Cannot Parse URL.")
			return
		}

		filepath := fmt.Sprintf("%s/", os.Getenv("TWITTER_DL_DIR"))
		filename := strings.Replace(u.Path, "media/", "", -1)

		err = DownloadFile(filepath, filename, imgURL, http.Header{})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error: Download Failed.")
			return
		}
	} else {
		filepath := fmt.Sprintf("%s/%s/", os.Getenv("TWITTER_DL_DIR"), twitterAPIData.Data[0].ID)
		for _, v := range twitterAPIData.Includes.Media {
			imgURL := v.URL
			u, err := url.Parse(imgURL)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error: Cannot Parse URL.")
				return
			}
			filename := strings.Replace(u.Path, "media/", "", -1)
			err = DownloadFile(filepath, filename, imgURL, http.Header{})
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error: Download Failed.")
				return
			}
		}
	}
	s.ChannelMessageSend(m.ChannelID, "Download Successfull.")
}
