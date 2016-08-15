package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/parnurzeal/gorequest"
)

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB

	hashSize = 10 * MB
)

type Subtitles struct {
	Status    string `xml:"status"`
	Subtitles struct {
		Content string `xml:"content"`
	} `xml:"subtitles"`
}

func getHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	rawContent := make([]byte, hashSize)
	_, err = f.Read(rawContent)
	if err != nil {
		return "", err
	}

	rawHash := md5.Sum(rawContent)
	hash := hex.EncodeToString(rawHash[:])

	return hash, nil
}

func parseResponse(response string) (string, error) {
	subtitles := &Subtitles{}
	if err := xml.Unmarshal([]byte(response), subtitles); err != nil {
		return "", err
	}
	if subtitles.Status != "success" || subtitles.Subtitles.Content == "" {
		return "", fmt.Errorf("Couldn't get subtitles: %s", subtitles.Status)
	}

	rawContent, err := base64.StdEncoding.DecodeString(subtitles.Subtitles.Content)
	if err != nil {
		return "", err
	}

	content := string(rawContent)

	return content, nil
}

func getSubtitles(path string) (string, error) {
	hash, err := getHash(path)
	if err != nil {
		return "", fmt.Errorf("Couldn't get hash from file: %v", err)
	}

	fmt.Println("Computed hash:", hash)

	// TODO: use SendStruct instead of Send
	resp, body, errs := gorequest.New().Post("http://napiprojekt.pl/api/api-napiprojekt3.php").
		Set("Content-Type", "application/x-www-form-urlencoded").
		Send("mode=1").
		Send("client=Napiprojekt").
		Send("client_ver=2.2.0.2399").
		Send("downloaded_subtitles_txt=1").
		Send("downloaded_subtitles_lang=PL").
		Send(fmt.Sprint("downloaded_subtitles_id=", hash)).
		Timeout(10 * time.Second).
		End()

	if (resp != nil && resp.StatusCode != 200) || body == "" {
		return "", fmt.Errorf("Couldn't get subtitles: %d, %s", resp.StatusCode, body)
	}
	if len(errs) > 0 {
		return "", fmt.Errorf("Couldn't get subtitles: %v", errs)
	}

	content, err := parseResponse(body)
	if err != nil {
		return "", fmt.Errorf("Couldn't parse response: %v", err)
	}

	return content, nil
}
