package main

import (
	"errors"
	"io"
	"net/http"
	"os"
)

// DownloadFile ファイルをダウンロードする汎用関数
func DownloadFile(filepath string, filename string, url string, headers http.Header) error {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header = headers

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		os.Mkdir(filepath, 0777)
	}

	out, err := os.Create(filepath + filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// RequestAjax APIリクエストする汎用関数
func RequestAjax(url string, headers http.Header) ([]byte, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header = headers

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("request failed: " + string(resp.StatusCode))
	}

	body, _ := io.ReadAll(resp.Body)
	return body, nil
}
