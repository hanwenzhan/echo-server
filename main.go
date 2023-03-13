// Copyright 2016 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {

		events, err := bot.ParseRequest(req)
		if err != nil {
			log.Println("server: ParseRequest: ", err)
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("server: could not read request body: %s\n", err)
		}
		log.Printf("server: request body: %s\n", reqBody)

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				source := event.Source
				log.Printf("server: UserID: %s\n", source.UserID) //, source.GroupID, source.RoomID)

				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					bot.BroadcastMessage(linebot.NewTextMessage("Hi " + source.UserID))

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Println("server: ", err)
					}

				case *linebot.StickerMessage:
					bot.BroadcastMessage(linebot.NewTextMessage("Hi Hi " + source.UserID))

					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Println("server: ", err)
					}
				}
			}
		}
	})

	//
	type Info struct {
		Message string `json:"message"`
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m := &Info{
			Message: "Hello World!",
		}
		b, err := json.Marshal(m)
		if err != nil {
			log.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	// This is just a sample code.
	// For actually use, you must support HTTPS by using `ListenAndServeTLS`, reverse proxy or etc.
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
