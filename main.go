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

		for _, event := range events {
			if event.Type != linebot.EventTypeMessage {
				continue
			}

			source := event.Source

			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				log.Printf("server: UserID: %s, Text: %s\n", source.UserID, message.Text)
				msg := fmt.Sprintf("your id is %s", source.UserID)
				//*
				_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(msg)).Do()
				if err != nil {
					log.Println("server: ReplyMessage: ", err)
				}
				//*/
				/*
					_, err = bot.PushMessage(source.UserID, linebot.NewTextMessage(msg)).Do()
					if err != nil {
						log.Println("server; PushMessage: ", err)
					}
				*/
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

	log.Printf("server version: %s \n", "0.0.2")
	//
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}

	log.Println("server exited.")
}
