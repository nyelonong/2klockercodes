package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type LockerCodes struct {
	Result struct {
		Data struct {
			AllLockerCodes struct {
				Edges []struct {
					Node struct {
						ID          string      `json:"id"`
						Title       string      `json:"title"`
						Version     int         `json:"version"`
						Slug        string      `json:"slug"`
						TweetID     string      `json:"tweetId"`
						DateAdded   string      `json:"dateAdded"`
						LockerCode  string      `json:"lockerCode"`
						Expiration  string      `json:"expiration"`
						LimitedTime interface{} `json:"limitedTime"`
						Reward      string      `json:"reward"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"allLockerCodes"`
		} `json:"data"`
	} `json:"result"`
}

func main() {
	telegramChatID := os.Getenv("TELEGRAM_ID")
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	lockerCodeHost := os.Getenv("LOCKERCODES_HOST")
	lockerCodePath := os.Getenv("LOCKERCODES_PATH")

	var version int
	flag.IntVar(&version, "version", 24, "NBA 2K Version")
	flag.Parse()

	if version == 0 {
		version = 24
	}

	jkt, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}

	mtn, err := time.LoadLocation("US/Mountain")
	if err != nil {
		panic(err)
	}

	link := fmt.Sprintf("%s/%d/%s", lockerCodeHost, version, lockerCodePath)

	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var lockerCodes LockerCodes
	if err := json.Unmarshal(body, &lockerCodes); err != nil {
		panic(err)
	}

	var hasData bool
	var msg []string
	msg = append(msg, fmt.Sprintf("NBA 2K%d", version))

	for _, edge := range lockerCodes.Result.Data.AllLockerCodes.Edges {
		code := edge.Node.LockerCode
		title := edge.Node.Title
		createdAt := edge.Node.DateAdded

		createdTime, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			panic(err)
		}

		now := time.Now().In(mtn)

		if now.Sub(createdTime).Hours() > 24 {
			continue
		}

		expire, err := time.Parse(time.RFC3339, edge.Node.Expiration)
		if err != nil {
			continue
		}

		if now.After(expire) {
			continue
		}

		msg = append(msg, fmt.Sprintf("Title: %s\nCode: %s\nCreated At: %s\nExpire At: %s\n\n", title, code, createdTime.In(jkt).String(), expire.In(jkt).String()))
		hasData = true
	}

	if !hasData {
		msg = append(msg, "No code today :(")
	}

	v := url.Values{}
	v.Set("chat_id", telegramChatID)
	v.Set("text", strings.Join(msg, "\n"))

	_, err = http.PostForm(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegramToken), v)
	if err != nil {
		panic(err)
	}
}
