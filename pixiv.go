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
	pixivURLPrefix   = "https://www.pixiv.net/"
	pixivAjaxBaseURL = "https://www.pixiv.net/ajax/illust/"
	RequestUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.60 Safari/537.36"
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

type pixivMangaAjaxRes struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Body    []struct {
		Urls struct {
			ThumbMini string `json:"thumb_mini"`
			Small     string `json:"small"`
			Regular   string `json:"regular"`
			Original  string `json:"original"`
		} `json:"urls"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"body"`
}

// DownloadPixivImg 貼り付けられたPixivのリンク先の画像をダウンロードするHandler
func DownloadPixivImg(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ボットユーザからのメッセージは何もしない
	if m.Author.ID == s.State.User.ID {
		return
	}

	// PixivのURLから始まらない場合は何もしない
	if !strings.HasPrefix(m.Content, pixivURLPrefix) {
		return
	}

	pixivURL := m.Content
	headers := http.Header{
		"referer":    []string{pixivURL},
		"user-agent": []string{RequestUserAgent},
	}

	// Pixivのイラスト用API URLを作成する
	pixivAjaxURL, err := makePixivAjaxURL(pixivURL)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error: Cannot Parse URL.")
		return
	}

	// Pixivのイラスト用API URLにリクエストする
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
		// マンガの場合
		// Pixivのマンガ用API URLを作成する
		mangaAjaxURL := "https://www.pixiv.net/ajax/illust/" + pixivAjaxData.Body.IllustID + "/pages"

		// Pixivのマンガ用API URLにリクエストする
		pixivMangaAjaxBytes, err := RequestAjax(mangaAjaxURL, headers)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error: Cannot Download Pixiv JSON.")
			return
		}

		// Pixivのマンガ用JSONをUnmarshalする
		var pixivMangaAjaxData pixivMangaAjaxRes
		json.Unmarshal(pixivMangaAjaxBytes, &pixivMangaAjaxData)

		// Pixiv側のエラーをハンドリングする
		if pixivMangaAjaxData.Error {
			s.ChannelMessageSend(m.ChannelID, "Error: Occurring Pixiv Error.")
			return
		}

		err = downloadPixivManga(pixivAjaxData, pixivMangaAjaxData, headers)
		if err != nil {
			fmt.Println(err)
			s.ChannelMessageSend(m.ChannelID, "Error: Failed download Pixiv images.")
			return
		}
	} else {
		// イラストの場合
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
	filepath := fmt.Sprintf("%s/", os.Getenv("PIXIV_DL_DIR"))
	filename := fmt.Sprintf("%s.png", pixivAjaxData.Body.IllustID)
	err := DownloadFile(filepath, filename, pixivAjaxData.Body.Urls.Original, headers)
	if err != nil {
		return err
	}
	return nil
}

func downloadPixivManga(pixivIllustAjaxData pixivIllustAjaxRes, pixivMangaAjaxData pixivMangaAjaxRes, headers http.Header) error {
	for i, v := range pixivMangaAjaxData.Body {
		filepath := fmt.Sprintf("%s/%s/", os.Getenv("PIXIV_DL_DIR"), pixivIllustAjaxData.Body.IllustID)
		filename := fmt.Sprintf("%d.png", i)
		err := DownloadFile(filepath, filename, v.Urls.Original, headers)
		if err != nil {
			return err
		}
	}
	return nil
}
