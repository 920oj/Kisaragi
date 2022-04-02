package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	pixivAjaxBaseURL = "https://www.pixiv.net/ajax/illust/"
)

type pixivIllustAjaxRes struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Body    struct {
		IllustID      string    `json:"illustId"`
		IllustTitle   string    `json:"illustTitle"`
		IllustComment string    `json:"illustComment"`
		ID            string    `json:"id"`
		Title         string    `json:"title"`
		Description   string    `json:"description"`
		IllustType    int       `json:"illustType"`
		CreateDate    time.Time `json:"createDate"`
		UploadDate    time.Time `json:"uploadDate"`
		Urls          struct {
			Mini     string `json:"mini"`
			Thumb    string `json:"thumb"`
			Small    string `json:"small"`
			Regular  string `json:"regular"`
			Original string `json:"original"`
		} `json:"urls"`
		UserID        string `json:"userId"`
		UserName      string `json:"userName"`
		UserAccount   string `json:"userAccount"`
		LikeData      bool   `json:"likeData"`
		Width         int    `json:"width"`
		Height        int    `json:"height"`
		PageCount     int    `json:"pageCount"`
		BookmarkCount int    `json:"bookmarkCount"`
		LikeCount     int    `json:"likeCount"`
		CommentCount  int    `json:"commentCount"`
		ResponseCount int    `json:"responseCount"`
		ViewCount     int    `json:"viewCount"`
	}
}

// DownloadPixivImg 貼り付けられたPixivのリンク先の画像をダウンロードするHandler
func DownloadPixivImg(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ボットユーザからのメッセージは何もしない
	if m.Author.ID == s.State.User.ID {
		return
	}

	// PixivのURLから始まらない場合は何もしない
	if !strings.HasPrefix(m.Content, "https://www.pixiv.net/") {
		return
	}

	pixivURL := m.Content
	headers := http.Header{
		"referer":    []string{pixivURL},
		"user-agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Safari/537.36"},
	}

	// Pixivのイラスト用APIリンクを作成する
	pixivAjaxURL, err := makePixivAjaxURL(pixivURL)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error: Cannot Parse URL.")
		return
	}

	// Pixivのイラスト用APIにリクエストする
	pixivAjaxBytes, err := RequestAjax(pixivAjaxURL, headers)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error: Cannot Download Pixiv JSON.")
		return
	}

	// PixivのJSONをUnmarshalする
	var pixivAjaxData pixivIllustAjaxRes
	json.Unmarshal(pixivAjaxBytes, &pixivAjaxData)

	// pixiv側のエラーをハンドリングする
	if pixivAjaxData.Error {
		s.ChannelMessageSend(m.ChannelID, "Error: Occurring Pixiv Error.")
		return
	}

	// イラストか漫画かの判定を行う
	if pixivAjaxData.Body.PageCount > 1 {
		// not implemented
		return
	} else {
		err := downloadPixivIllust(pixivAjaxData, headers)
		if err != nil {
			fmt.Println(err)
			s.ChannelMessageSend(m.ChannelID, "Error: Failed download Pixiv image.")
			return
		}
	}
	s.ChannelMessageSend(m.ChannelID, "Download Successfull.")
}

func makePixivAjaxURL(pixivURL string) (string, error) {
	// URLをパースして、イラストIDを取得する
	u, err := url.Parse(pixivURL)
	if err != nil {
		return "", err
	}
	illustId := strings.Replace(u.Path, "/artworks/", "", -1)

	// AJAXのリンクを生成する
	pixivAjaxURL := pixivAjaxBaseURL + illustId
	return pixivAjaxURL, nil
}

func downloadPixivIllust(pixivAjaxData pixivIllustAjaxRes, headers http.Header) error {
	err := DownloadFile(os.Getenv("PIXIV_DL_DIR")+"/"+pixivAjaxData.Body.IllustID+".png", pixivAjaxData.Body.Urls.Original, headers)
	if err != nil {
		return err
	}
	return nil
}
