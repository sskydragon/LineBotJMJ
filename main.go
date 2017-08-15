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
	"regexp"
	"strconv"
	
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
		case (strings.Contains(msg,"!指令")):
			replyMsg = "測試 / 同好會社團 / 日本麻將介紹網站 / \n"+
						"日麻行事曆 / 過去的活動 / [IORMC|WRC|雀鳳盃|般若盃]資訊\n"+
						"覺得有漏可以tag龍哥幫忙加, 還有一些小彩蛋就不說了喵~"
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
						"今年第四屆報名時間已經過囉 下次請早喔喵QwQ"
		case (strings.Contains(msg,"雀鳳盃資訊")):
			replyMsg = "雀鳳盃是目前臺灣北部最大的例行比賽 通常在三月\n"+
						"第四屆/2017年辦法請參考 http://goo.gl/SB4yth\n"+
						"好期待明年的比賽呢喵~ (滾來滾去)"
		case (strings.Contains(msg,"日麻行事曆")):
			replyMsg = "在這在這~~ http://goo.gl/fqwswg"
		case (strings.Contains(msg,"過去的活動")):
			replyMsg = "喵知道過去的比賽活動有這些~！\n"+
						"https://goo.gl/KH00SO\n" +
						"還想提供些什麼的話也請讓喵知道喔喵~ "
		case (strings.Contains(msg,"萌新")):
			replyMsg = "萌新是在說我嗎喵~ (探頭"
		case (strings.Contains(msg,"池田銀行")):
			replyMsg = "點數太多嗎？歡迎存點數進來悠喵OwO"
		case (strings.Contains(msg,"何切")):
			words := strings.Fields(msg)
			re := regexp.MustCompile("(([0-9]+[MmPpSs])|([1-7]+[Zz]))+")
			result := "";
			reply := "";
			status := 0;
			for i := 0; i < len(words); i++ {
				result = re.FindString(words[i]);
				numAmount := 0;

				for j := 0 ; j < len(result) ; j++ {
					_, err := strconv.ParseFloat(string(result[j]), 64)			
					if(err == nil) {numAmount++}
				}
				if numAmount > 0 {
					if numAmount % 3 != 2 {
					 status = -1
					 reply="這樣好像不能拿去問天鳳姬呢喵~\n"+"手牌必須是0~9接花色mpsz的組合\n"+
					 "(0是赤 m萬p筒s索z字 字牌只有1-7)\n"+"而且丟完牌必須是3n+1張才不會出錯唷~"
					} else {
						status = 1
						reply = "天鳳姬是這樣說的呢喵~\n"+"http://tenhou.net/2/?q="+result+"\n"
						break
					}
				}
			}
			if(status != 0) {replyMsg = reply}
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
				replyMsg := determineReply(message.Text)
				if replyMsg != "" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMsg)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}
