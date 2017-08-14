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
		case (strings.Contains(msg,"婊池田")),(strings.Contains(msg,"打爆池田")):
			replyMsg = "不要欺負池田喵好嗎 QQ"
		case (strings.Contains(msg,"龍哥救我")):
			replyMsg = "需要幫忙嗎喵~？"
		case (strings.Contains(msg,"IORMC資訊")):
			replyMsg =  "IORMC是韓國辦的國際交流賽\n" +
						"2017的預選資訊在 https://goo.gl/2XJyYw\n" +
						"2016的比賽結果在 https://goo.gl/jatIHN"
		case (strings.Contains(msg,"WRC資訊")):
			replyMsg =  "WRC是三年一次的世界日麻大賽(暫譯)\n"+
						"2017.10.4-8在拉斯維加斯, 網站在 http://www.wrc2017vegas.com/\n"+
						"預計2020年在歐洲、2023年在日本舉行"
		case (strings.Contains(msg,"同好會社團")):
			replyMsg = "我們的社團在這喔喵つ https://www.facebook.com/groups/twjmj/"
		case (strings.Contains(msg,"般若盃資訊")):
			replyMsg = "般若盃是目前臺灣南部最大的例行比賽 通常在十月\n"+
						"簡章請參考 https://goo.gl/XjHCfW\n"+
						"今年報名時間已經過囉 下次請早喔喵QwQ"
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
