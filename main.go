// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func determineReply(msg string) string{
	var replyMsg string = ""
	switch {
		case (strings.Contains(msg,"測試")):
			replyMsg = "在測試啥呢喵~"
		case (strings.Contains(msg,"日本麻將介紹網站")):
			replyMsg = "介紹網站在 http://jmj.tw 喵~"
		case (strings.Contains(msg,"婊池田")):
		case (strings.Contains(msg,"打爆池田")):
			replyMsg = "不要欺負池田喵好嗎 QQ"
		default:
	}
	return replyMsg
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(determineReply(message.Text))).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
