package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

type LineMessage struct {
	Destination string `json:"destination"`
	Events      []struct {
		ReplyToken string `json:"replyToken"`
		Type       string `json:"type"`
		Timestamp  int64  `json:"timestamp"`
		Source     struct {
			Type   string `json:"type"`
			UserID string `json:"userId"`
		} `json:"source"`
		Message struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"events"`
}

type ReplyMessage struct {
	ReplyToken string `json:"replyToken"`
	Messages   []Text `json:"messages"`
}

type Text struct {
	Type      string `json:"type"`
	Text      string `json:"text"`
	PackageId string `json:"packageId"`
	StickerId string `json:"stickerId"`
}

type ProFile struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

var ChannelToken = os.Getenv("CHANNEL_TOKEN")

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.POST("/webhook", func(c echo.Context) error {

		Line := new(LineMessage)
		if err := c.Bind(Line); err != nil {
			log.Println("err", err)
			return c.String(http.StatusOK, "error")
		}
		fullname := getProfile(Line.Events[0].Source.UserID)

		str1 := "hi hello"
		text := Text{}
		if strings.Contains(str1, strings.ToLower(Line.Events[0].Message.Text)) {
			text = Text{
				Type: "text",
				Text: Line.Events[0].Message.Text + " " + fullname,
			}
			message := ReplyMessage{
				ReplyToken: Line.Events[0].ReplyToken,
				Messages: []Text{
					text,
				},
			}
			replyMessageLine(message)
		}

		str2 := "wow"
		text2 := Text{}
		if strings.Contains(str2, strings.ToLower(Line.Events[0].Message.Text)) {
			text2 = Text{
				Type:      "sticker",
				PackageId: "11537",
				StickerId: "52002734",
			}
			message := ReplyMessage{
				ReplyToken: Line.Events[0].ReplyToken,
				Messages: []Text{
					text2,
				},
			}
			replyMessageLine(message)
		}

		log.Println("%% message success")
		return c.String(http.StatusOK, "ok")

	})

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	fmt.Println("addr: ", addr)
	log.Println("starting...")
	e.Logger.Fatal(e.Start(addr))
}

func replyMessageLine(Message ReplyMessage) error {
	value, _ := json.Marshal(Message)

	url := "https://api.line.me/v2/bot/message/reply"

	var jsonStr = []byte(value)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+ChannelToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

	return err
}

func getProfile(userId string) string {

	url := "https://api.line.me/v2/bot/profile/" + userId

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+ChannelToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var profile ProFile
	if err := json.Unmarshal(body, &profile); err != nil {
		log.Println("%% err \n")
	}
	log.Println(profile.DisplayName)
	return profile.DisplayName

}
