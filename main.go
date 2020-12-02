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

/* required environment vars
ChannelAccessToken
ChannelSecret
AdminLineIDList   	(Line UserID)
SupportedGroups		(Line GroupID)
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"os"
	"regexp"
	"strconv"
	"time"

//	"io/ioutil"
//	"bytes"

	"github.com/line/line-bot-sdk-go/linebot"
//	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
)
//global cd會造成不同對話群組間互相影響的問題(再研究)
var bot *linebot.Client
var cdTest = 60*time.Second
var cdNewbie = 60*time.Second
var cdBank = 60*time.Second
var cdTeachMe = 60*time.Second
var cdLobby = 60*time.Second
var cdL1120 = 10*time.Second
var cdBullyCat = 60*time.Second
var cdSaveMe = 60*time.Second
var cdGiveUp = 60*time.Second
var cdWhatCutHelp = 60*time.Second
var cdSlides = 60*time.Second
var lastTest = time.Now().Add(-cdTest)
var lastNewbie = time.Now().Add(-cdNewbie)
var lastBank = time.Now().Add(-cdBank)
var lastTeachMe = time.Now().Add(-cdTeachMe)
var lastLobby = time.Now().Add(-cdLobby)
var lastL1120 = time.Now().Add(-cdL1120)
var lastBullyCat = time.Now().Add(-cdBullyCat)
var lastSaveMe = time.Now().Add(-cdSaveMe)
var lastGiveUp = time.Now().Add(-cdGiveUp)
var lastWhatCutHelp = time.Now().Add(-cdWhatCutHelp)
var lastSlides = time.Now().Add(-cdSlides)

var _001_help = regexp.MustCompile("!((指|命)令|(幫|帮|幇|鞤)助|菜(单|單|単)|功能|用法|使用方法|(説|说|說)明)")
var _002_website = regexp.MustCompile("(!(介紹)?網站)|((日麻|麻將).*介紹.*網站)")
var _003_bullycat = regexp.MustCompile("(婊|打).*池田")
var _004_saveme = regexp.MustCompile("救我")
var _005_iormc = regexp.MustCompile("(!iormc)|iormc.*資訊")
var _006_wrc = regexp.MustCompile("(!wrc)|wrc.*資訊")
var _007_fbgroup = regexp.MustCompile("(!(社團|同好會))|(社團|同好會).*(在哪|在那)")
var _008_banro = regexp.MustCompile("(!般若)|般若.*資訊")
var _009_janfun = regexp.MustCompile("(!雀鳳)|雀鳳.*資訊")
var _010_calendar = regexp.MustCompile("(!行事曆)|日麻行事曆")
var _011_pastactivity = regexp.MustCompile("(!(歷史|過去)?活動)|過去.*活動")
var _012_faq = regexp.MustCompile("(!(常見|新手)問題)|(新手|常見(的)?)問題")
var _013_tanyao = regexp.MustCompile("!斷(么|公)(九)?")
var _014_tsumo = regexp.MustCompile("!(自摸|門(前)?清)")
var _015_notexist = regexp.MustCompile("!(全求|花龍|一色四順|四歸一|組合龍)")
var _016_onlywaitone = regexp.MustCompile("!(單釣|單騎|獨聽|嵌張|崁張|邊張)")
var _017_iitsu = regexp.MustCompile("!(一氣|一通|一條龍)")
var _018_nagashimangan = regexp.MustCompile("!(流滿|流局滿貫|流し滿貫)")
var _019_chonbo = regexp.MustCompile("!(醬爆|錯和|チョンボ)")
var _020_kyotaku = regexp.MustCompile("!(供託|託供|供托|托供|きょうたく)")
var _021_2to8 = regexp.MustCompile("!(大車輪|大數鄰|大竹林|大萬華|大数隣)")
var _022_1to7or3to9 = regexp.MustCompile("!(小車輪|小數鄰|小竹林|小萬華|小数隣)")
var _023_following = regexp.MustCompile("!(尾行|真似滿|まねまん)")
var _024_usage = regexp.MustCompile("使用說明|用法")
var _025_wheretoplay = regexp.MustCompile("(哪|地方).*(玩|下)")
var _026_wheretoplayonmobile = regexp.MustCompile("(哪|地方).*(手機|app)")
var _027_ikedaok = regexp.MustCompile("摸摸池田(的)?(肚子|肚肚|頭|耳朵|尾巴|額頭|下巴)")
var _028_ikedang = regexp.MustCompile("戳|揍|扁|k|胸|屁|內|陰|婊|打|揍|胖次|歐派|腿|雞|懶|p|j|gg|bra")
var _029_interact = regexp.MustCompile("(摸摸|抱抱|拉拉)池田")
var _030_askDragon = regexp.MustCompile("!(天空龍|龍龍|蘿莉龍)")
var _031_askStarflyx = regexp.MustCompile("!星野")
var _032_askBingo = regexp.MustCompile("!賓果")
var _033_askTaiwanco = regexp.MustCompile("!(少年|沈皇|沈凰)")

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func reply1msg1video(replyToken string , msg string , videoSrc string , picSrc string ) {
	var err error
	if _, err = bot.ReplyMessage(replyToken, linebot.NewTextMessage(msg), linebot.NewVideoMessage(videoSrc, picSrc)).Do(); err != nil {
		log.Print(err)
	}

}

func reply1msg1pic(replyToken string, msg string, picSrc string) {
	var err error
	if _, err = bot.ReplyMessage(replyToken, linebot.NewTextMessage(msg), linebot.NewImageMessage(picSrc, picSrc)).Do(); err != nil {
		log.Print(err)
	}

}

func isAdmin(msg string) bool {
	s := os.Getenv("AdminLineIDList")
	if(strings.Contains(s,msg)) {
		return true
	}
	return false
}

func isSupportedGroup(msg string) bool {
	s := os.Getenv("SupportedGroups")
	if(strings.Contains(s,msg)) {
		return true
	}
	return false
}

func isExcludedGroup(msg string) bool {
	s := os.Getenv("ExcludedGroups")
	if(strings.Contains(s,msg)) {
		return true
	}
	return false
}

func teachMe(msg string) bool {
	if ((strings.Contains(msg,"日麻") || strings.Contains(msg,"麻將")) && strings.Contains(msg,"教學") && (strings.Contains(msg,"嗎") || strings.Contains(msg,"哪"))) {
		return true
	}
	if ((strings.Contains(msg,"日麻") || strings.Contains(msg,"麻將")) && (strings.Contains(msg,"想") || strings.Contains(msg,"要")) && strings.Contains(msg,"學") && strings.Contains(msg,"我")) {
		return true
	}
	if ((strings.Contains(msg,"日麻") || strings.Contains(msg,"麻將")) && strings.Contains(msg,"教我")){
		return true
	}
	if (strings.Contains(msg,"!教學")) {
		return true
	}
	return false
}

func askingLobby(msg string) bool {
/* 時效已過
	if(strings.Contains(msg,"!大會")) {
		return true
	}
	if(strings.Contains(msg,"大會室")&&strings.Contains(msg,"在哪")) {
		return true
	}
	if(strings.Contains(msg,"大會室")&&(strings.Contains(msg,"連結") || strings.Contains(msg,"網址")) &&(!strings.Contains(msg,".net") && !strings.Contains(msg,"http"))) {
		return true
	}
*/
	return false
}

func askingL1120(msg string) bool {
	if(strings.Contains(msg,"!個室")) {
		return true
	}
	if(strings.Contains(msg,"個室")&&strings.Contains(msg,"在哪")) {
		return true
	}
	if(strings.Contains(msg,"個室")&&(strings.Contains(msg,"連結") || strings.Contains(msg,"網址") || strings.Contains(msg,"位置")) &&(!strings.Contains(msg,".net") && !strings.Contains(msg,"http"))) {
		return true
	}

	return false
}

func askingNTUSlides(msg string) bool {
	if(strings.Contains(msg,"!台大講義") || strings.Contains(msg,"!臺大講義") || strings.Contains(msg,"!社課講義") || strings.Contains(msg,"!賓果講義") || strings.Contains(msg,"!上課錄影") ){
		return true
	}
	if( (strings.Contains(msg,"台大") || strings.Contains(msg,"臺大") || strings.Contains(msg,"賓果")) && (strings.Contains(msg,"講義") || strings.Contains(msg,"錄影") || strings.Contains(msg,"教材") || strings.Contains(msg,"上課內容")) && (strings.Contains(msg,"在哪") || strings.Contains(msg,"位置") || strings.Contains(msg,"給一下") || strings.Contains(msg,"哪找") || strings.Contains(msg,"？")) || strings.Contains(msg,"哪裡找得到") ) {
		return true
	}
	return false
}

func appendSkyInfo(msg string) string {
	msg += "天空龍的頻道上可以找到同好們的影片分享\n" +
			"內容已由淺至深整理，找教學先來這吧！" +
			"https://www.youtube.com/c/sskytw"
	return msg
}

func appendNTUSlidesInfo(msg string) string {
	msg += "想在理論上鑽研的話，可以研讀賓果寫的講義\n" +
			"上學期 https://goo.gl/bFBy9w\n" +
			"下學期 https://goo.gl/E9rirQ\n"
//			+"社課錄影 https://goo.gl/sYS6Vd"
	return msg
}

func appendStarflyxInfo(msg string) string {
	msg += "星野的教學影片 - 適合任何人觀看喔~\nhttps://goo.gl/Yjeixi"
	return msg
}

func appendTaiwancoInfo(msg string) string {
	msg += "少年與沈欸的天鳳觀戰解析 - 對特桌以上玩家較有幫助\nhttps://goo.gl/5PX5VH"
	return msg
}

func determineReply(msg string, groupSupported bool, groupExcluded bool) string{

/*	groupSupported是用來限制哪些功能是給日麻相關群組用的
	groupExcluded則是用來限制太過閒聊的功能在要專心討論的群組裡不要開放
*/
	var replyMsg string = ""
	t := time.Now()
	switch {
		case (_001_help.MatchString(msg)):
			replyMsg = "指令都用!開頭, 目前只能使用繁體詢問\n"+ 		   "我會告訴你一些我知道的東西喔喵~\n"+
						"學習 | !教學 !何切 !討論 !常見問題 ![役種名稱]\n"+
						"遊玩 | !天鳳 !雀姬 !雀魂 !個室 !戰績網 !門清狂 !BAMBOO\n"+
						"社群 | !社團 !網站\n"+
						"活動 | !IORMC !WRC !雀鳳 !般若 !行事曆 !歷史活動\n"+
						"其他 | !黑白棋教學 !氣泡教學 !方塊教學 (!?)\n"+
						"喵想知道大家還對什麼有興趣！快叫天空龍教我吧~~(翻滾)"
		case (_002_website.MatchString(msg)):
			replyMsg = "日本麻將介紹網站在 http://jmj.tw 喔喵~"
		case (t.Sub(lastBullyCat) > cdBullyCat && _003_bullycat.MatchString(msg)):
			lastBullyCat = t
			replyMsg = "不要欺負池田喵好嗎 QQ"
		case (t.Sub(lastSaveMe) > cdSaveMe && _004_saveme.MatchString(msg)):
			lastSaveMe = t
			replyMsg = "需要幫忙嗎喵~？"
		case (_005_iormc.MatchString(msg)):
			replyMsg =  "IORMC是韓國辦的國際交流賽\n" +
						"歷年紀錄請參考：\n"+
						"2018 https://goo.gl/YdzR5A\n"+
						"2017 https://goo.gl/VXvK3W\n"+
						"2016 https://goo.gl/jatIHN"
		case (_006_wrc.MatchString(msg)):
			replyMsg =  "WRC是三年一次的世界日麻大賽(暫譯)\n"+
						"2017.10.4-8在拉斯維加斯, 網站在 http://www.wrc2017vegas.com/\n"+
						"預計2020年在歐洲、2023年在日本舉行"
		case (_007_fbgroup.MatchString(msg)):
			replyMsg = "我們的社團在這喔喵つ https://www.facebook.com/groups/twjmj/"
		case (t.Sub(lastLobby) > cdLobby && askingLobby(msg)):
			lastLobby = t
			replyMsg = "限時開放的IORMC練習大會室在這喔喵~\n"+
			"http://tenhou.net/0/?C85193656"
		case (t.Sub(lastL1120) > cdL1120 && askingL1120(msg)):
			lastL1120 = t
			replyMsg = "平常用來交流的個室在這喔喵~\n"+
			"http://tenhou.net/0/?L1120"
		case (_008_banro.MatchString(msg)):
			replyMsg = "般若盃是目前臺灣南部最大的例行比賽 通常在十月\n"+
						"報名超額時會依天鳳段位及牌譜內容決定參與者\n"+
						"把天鳳段位打高一點比較容易報上喔~"
		case (_009_janfun.MatchString(msg)):
			replyMsg = "雀鳳盃是目前臺灣北部最大的例行比賽 通常在三月\n"+
						"相關資訊能在淡大日麻社社團找到 https://goo.gl/9FFvvn\n"
		case (_010_calendar.MatchString(msg)):
			replyMsg = "在這在這~~ http://goo.gl/fqwswg"
		case (_011_pastactivity.MatchString(msg)):
			replyMsg = "喵知道過去的比賽活動有這些~！\n"+
						"https://goo.gl/KH00SO\n" +
						"還想提供些什麼的話也請讓喵知道喔喵~ "
		case (_012_faq.MatchString(msg)):
			replyMsg = "※和牌必須要有役(寶牌以外)\n"+
			"※不能振聽(聽的牌中不能有自己打出過的牌。\n"+
			"  若別人打出的牌能讓自家湊成和牌型、但自家不能和或選擇不和時會暫時振聽,\n"+
			"  非立直狀態下要到自己摸牌時才解除；立直後則是永久振聽)\n"+
			"※高取牌(和了時必須先把牌型拆開擺好, 且只計算最有利的一種)\n"+
			"※只計上位役(二盃口必含一盃口, 所以有二盃的情況不計一盃)\n"+
			"※鳴牌降飜(部分役種在有叫牌的情況下價值會下降或不計)\n"+
			"※役滿只能和役滿複合(有役滿的情況下不計一般役種)\n\n"+

			"各役種常見問題 請用「!役種名稱」查詢\n"
		case (_013_tanyao.MatchString(msg)):
			replyMsg = "斷么只能有2~8的數字牌"
		case (strings.Contains(msg,"!一盃")):
			replyMsg = "一盃口必須門清才計算"
		case (strings.Contains(msg,"!二盃")):
			replyMsg = "二盃口必須門清才計算\n"+
						"因為拆解型不同, 二盃口不計七對子"
		case (strings.Contains(msg,"!七對")):
			replyMsg = "七對必須有七組「不同」的對子, 符數是固定25符\n"+
						"因為拆解型不同, 所以採用七對時不計算一般型的役種"
		case (strings.Contains(msg,"!平和") || strings.Contains(msg,"!平胡")):
			replyMsg = "平和必須門前清、聽兩面且牌型沒有符\n"+
						"(整手牌都是123這樣的順面組合 而且雀頭不是役牌)"
		case (strings.Contains(msg,"!門清狂")):
			replyMsg = "門清狂是熟悉清一色的好遊戲~！\n"+
						"http://hinakin.main.jp/mckonweb/index.htm"
		case (strings.Contains(msg,"!bamboo")):
			replyMsg = "BAMBOO是清一色對戰麻雀, 破關有隱藏模式可以期待唷~\n"+
						"http://www.gamedesign.jp/flash/bamboo/bamboo.html"
		case (_014_tsumo.MatchString(msg)):
			replyMsg = "日麻自摸和門清都沒有役,\n"+
					   "只有在門前清的情況下自摸才有一飜"
		case (strings.Contains(msg,"!寶牌") || strings.Contains(msg,"!ドラ")):
			replyMsg = "寶牌是「寶牌指示牌」的下一張, 且不能當起和飜\n"+
						"1→2→..→9→1, 東→南→西→北→東 白→發→中→白\n"+
						"裡寶牌只有立直且和出的人才計算"
		case (_015_notexist.MatchString(msg)):
			replyMsg = "日麻沒有這種東西...."
		case (_016_onlywaitone.MatchString(msg)):
			replyMsg = "日麻只聽一張的情況並沒有飜, 只在符數上會有一點點加成"
		case (strings.Contains(msg,"!立直")):
			replyMsg = "立直必須滿足三個條件：\n"+
						"(1)門前清聽牌(不能有牌從別人叫來)\n"+
						"(2)能夠支付立直供託(通常是1000點)\n"+
						"(3)還預期有牌能摸(至少剩四張可摸牌)"
		case (strings.Contains(msg,"!一發")):
			replyMsg = "立直後到自己下次捨牌前和了, 且過程中沒有人鳴牌才計算"
		case (strings.Contains(msg,"!役牌")):
			replyMsg = "場風/自風/三元(白發中)有三張以上才算(看錯場風自風了嗎?)"
		case (_017_iitsu.MatchString(msg)):
			replyMsg = "和了時的拆解 必須能拆出123 456 789三組同色的順子 才有一氣"
		case (_018_nagashimangan.MatchString(msg)):
			replyMsg = "捨牌全是19字、且沒有被人叫走才算。\n"+
						"在天鳳和流局聽牌狀態分開計算點數增減。不一定有玩的規則。"
		case (strings.Contains(msg,"!四暗")):
			replyMsg = "四暗刻的四組刻子都要自力摸進來才算\n"+
						"最後一張榮和是來自別人的話, 要算明刻, 沒有四暗"
		case (strings.Contains(msg,"!綠一")):
			replyMsg = "(地方役)(常見) 綠一只能有23468索和發。通常不需要有發。"
		case (strings.Contains(msg,"!人和")):
			replyMsg = "(地方役) 自己摸到牌前、無人鳴牌的情況下榮和(天鳳不採用)"
		case (strings.Contains(msg,"!搶槓") || strings.Contains(msg,"!槍槓")):
			replyMsg = "加槓的那張牌剛好是和了牌。因為槓牌失敗了, 所以不會開新寶牌。\n"+
						"部分規則下允許國士無雙搶暗槓, 除此之外沒有例外(一般不能搶暗槓)"
		case (strings.Contains(msg,"!加槓")):
			replyMsg = "已經碰出一組牌後, 再摸到相同的牌可以「加槓」, 從牌山後方補牌\n"+
						"開新寶牌的時機一般有兩種, 一種是補牌前開, 一種是捨牌/再槓時同時開"
		case (strings.Contains(msg,"!大明槓")):
			replyMsg = "手中拿出三張一樣的牌, 去槓別人丟出的第四張牌\n"+
						"開新寶牌的時機一般有兩種, 一種是補牌前開, 一種是捨牌/再槓時同時開"
		case (strings.Contains(msg,"!役種")):
			replyMsg = "一般役種說明請參考 http://jmj.tw/extra/yakutable.pdf\n"	 +
						"除此之外還有些地方役, 不見得每個地方都玩, 要先問規則喔~"
		case (strings.Contains(msg,"!段位")):
			replyMsg = "各平台都有自己的段位系統, 天鳳的較有參考價值\n" +
						"天鳳的段級位制請參考 http://tenhou.net/man/#DAN"
		case (_019_chonbo.MatchString(msg)):
			replyMsg = "錯和/チョンボ/醬爆\n"+
						"通常是指在不能和牌的情況下和牌、或立直卻沒聽等嚴重影響牌局的錯行為\n" +
						"罰則視規定不同, 一般有罰自摸滿貫 或是結束後扣大量分數的兩種做法"
		case (strings.Contains(msg,"!相公")):
			replyMsg = "多牌少牌摸錯牌食替或其他不非常嚴重但足以影響牌局的行為\n"+
						"會視為相公, 一般懲罰是無法鳴牌、和牌, 流局時也不算聽牌"
		case (strings.Contains(msg,"!食替")):
			replyMsg = "吃牌後立刻扔直接相關的牌\n"+
						"例如手中有3456m，用45m吃了6m後立刻扔3m或6m (都和45m直接相關)"
		case (strings.Contains(msg,"!送槓")):
			replyMsg = "立直後槓牌破壞牌型的行為。\n"+
						"一般立直後的槓牌不能影響聽牌種類,\n"+
						"也有更嚴格的規定如不能影響役種或讓拆解方式消失等"
		case (strings.Contains(msg,"!叫牌換牌") || strings.Contains(msg,"!叫換牌")):
			replyMsg = "實際行為與出聲不同, 例如喊碰後改用吃的方式拿牌"
		case (strings.Contains(msg,"!叫牌修正")):
			replyMsg = "宣告鳴牌後, 因為無法鳴下或突然不想鳴等原因取消"
		case (_020_kyotaku.MatchString(msg)):
			replyMsg = "「供託」：放在場上的千點棒, 包括立直時付出的千點、和輕微犯規時罰的千點\n"+
						"只有「供託」是正確的寫法喔！"
		case (strings.Contains(msg,"!積棒") || strings.Contains(msg,"!本場棒")):
			replyMsg = "用來計算連莊次數的百點棒。\n"+
						"每次正常流局或親家連莊時要多放一根, 子家和了時歸零"

		case (strings.Contains(msg,"!三色同刻") || strings.Contains(msg,"!三同刻")):
			replyMsg = "三種顏色都有相同數字的刻子(刻子指三張相同牌, 如222m)。\n"+
						"一般講三色是指「三色同順」"
		case (strings.Contains(msg,"!刻子") || strings.Contains(msg,"!暗刻") || strings.Contains(msg,"!明刻")):
			replyMsg = "刻子指一組三張相同的牌, 例如222m。\n"+
						"三張都是自己摸進來的話是「暗刻」, 有任一張從對手來要算「明刻」"
		case (strings.Contains(msg,"!順子")):
			replyMsg = "順子指一組三張同色相連但不同數字的牌, 例如123m"
		case (strings.Contains(msg,"!三連刻")):
			replyMsg = "(地方役) 和了型中 包含同一色三組數字相連的刻子 例如111m 222m 333m"
		case (strings.Contains(msg,"!四連刻")):
			replyMsg = "(地方役) 和了型中 包含同一色四組數字相連的刻子 例如111m 222m 333m 444m"
		case (strings.Contains(msg,"!地方役")):
			replyMsg = "有些地方玩、而有些地方不玩的規則"
		case (strings.Contains(msg,"!百萬石")):
			replyMsg = "(地方役) 清一色且數牌加起來總和有到100, 以役滿計算\n"+
						"若總和恰為100, 稱為「加賀百万石」或「純正百万石」, 通常算雙役滿"
		case (strings.Contains(msg,"!開牌立直") || strings.Contains(msg,"!開立")):
			replyMsg = "(地方役) 立直時打開全部或部分的手牌(依規定不同), 算兩飜\n"+
						"若非立直狀態下銃了開立直的人, 以役滿計算"
		case (strings.Contains(msg,"!八連莊") || strings.Contains(msg,"!阻止八連莊") || strings.Contains(msg,"!破回八連莊") ):
			replyMsg = "(地方役)(罕見) 親家「連續和了」八次時, 第八次(含以後)不論牌型皆以役滿計算\n"+
						"只有和了才算, 中間有流局要重新計算。\n"+
						"有些規則會把阻止親家的第八次和了也當作役滿"
		case (strings.Contains(msg,"!一色三步高")):
			replyMsg = "(這是國標的番種, 日麻沒這種東西) 同色以差距1或2遞增的三組順子\n"+
						"例如123m 234m 345m或是123m 345m 567m"
		case (strings.Contains(msg,"!三色三步高")):
			replyMsg = "(這是國標的番種, 日麻沒這種東西) 三色分別有一組順子, 以差距1或2遞增\n"+
						"例如123m 234p 345s或是123m 345p 567s"
		case (strings.Contains(msg,"!一色三節高")):
			replyMsg = "(這是國標的番種, 日麻沒這種東西) 相當於日麻的三連刻, 同色三組相連的刻子\n"+
						"例如111m 222m 333m"
		case (strings.Contains(msg,"!三色三節高")):
			replyMsg = "(這是國標的番種, 日麻沒這種東西) 三種顏色分別有三組數字差1的刻子\n"+
						"例如111m 222p 333s"
		case (strings.Contains(msg,"!清老頭")):
			replyMsg = "和牌時手牌全由19數牌組成。役滿。(有字的話是混老頭)"
		case (strings.Contains(msg,"!混老頭")):
			replyMsg = "和牌時手牌全由19數牌和字牌組成。(沒有字的話是清老頭)"
		case (strings.Contains(msg,"!金門橋") || strings.Contains(msg,"!金門大橋")):
			replyMsg = "(地方役)(罕見) 和了時包含同一色的123 345 567 789四組順子。Local役滿。\n"+
						"沒玩這個規則的話, 可能會連和牌都不行呢...."
		case (strings.Contains(msg,"!黑一色")):
			replyMsg = "(地方役)(罕見) 和了時只有248筒和東南西北風牌。Local役滿。"
		case (strings.Contains(msg,"!紅孔雀")):
			replyMsg = "(地方役)(罕見) 和了時只有1579索和中。Local役滿。"
		case (strings.Contains(msg,"!大三索")):
			replyMsg = "參考 https://zh.moegirl.org/zh-tw/三索必须死\n"+
						"科學麻將不考慮這種超能力麻將情節的....(汗"
		case (strings.Contains(msg,"!大七星")):
			replyMsg = "(地方役) 東東南南西西北北白白發發中中 的七對子。Local雙役滿。"
		case (strings.Contains(msg,"!石上三年") || strings.Contains(msg,"!鐵杵成針") || strings.Contains(msg,"!石の上にも三年")):
			replyMsg = "(地方役)(罕見) 雙立直+海底撈月同時成立。Local役滿。"
		case (strings.Contains(msg,"!超立直")):
			replyMsg = "(地方役)(罕見) 支付五千點立直, 和了時寶牌指示牌前後一張都算寶牌。"
		case (_023_following.MatchString(msg)):
			replyMsg = "(地方役)(罕見) 無人鳴牌的情況下, 和風位在前的玩家捨牌完全相同\n"+
						"一般至少要跟出四張才能和, 跟出一張算一飜"
		case (_021_2to8.MatchString(msg)):
			replyMsg = "(地方役) 2-8各兩張的門前清清一色。\n"+
						"根據花色不同, 而有大數鄰or大萬華(萬)/大車輪(筒)/大竹林(索)的稱呼。\n"+
						"Local役滿。"
		case (_022_1to7or3to9.MatchString(msg)):
			replyMsg = "(地方役) 1-7或3-9各兩張的門前清清一色。(偏一邊)\n"+
						"根據花色不同, 而有小數鄰or小萬華(萬)/小車輪(筒)/小竹林(索)的稱呼。\n"+
						"Local役滿。"
		case (strings.Contains(msg,"!東北新幹線")):
			replyMsg = "(地方役)(罕見) 含有東、北及一氣通貫的混一色\n"+
						"一般必須門前清, 視規則可能只承認筒子和/或索子。Local役滿。"
		case (strings.Contains(msg,"!十三不搭") || strings.Contains(msg,"!十三不塔")):
			replyMsg = "(地方役) 起手摸進牌後除了一對以外無法形成任何搭子且先前沒人鳴牌。Local役滿。"
		case (strings.Contains(msg,"!十四不靠") || strings.Contains(msg,"!十四無靠")):
			replyMsg = "(地方役) 起手摸進牌後無法形成任何搭子且先前沒人鳴牌。Local役滿。"
		case (strings.Contains(msg,"!圓周率")):
			replyMsg = "(地方役) 圓周率役滿：\n圓周率(π) 小數點以下前14位.14159265358979\n"+
			"重新排列後正好是 11 234 555 678 999 的和牌形，\n所以稱圓周率。Local役滿。"
		case (strings.Contains(msg,"!燕返") || strings.Contains(msg,"!燕返し")):
			replyMsg = "(地方役)(罕見) 和別人的立直宣言牌。一飜。\n"+
						"(一般講燕返し不是指地方役, 而是從牌山換回整副手牌的作弊方法)"


		case (strings.Contains(msg,"!螢返") || strings.Contains(msg,"!蛍返")):
			replyMsg = "一種帥氣的倒牌方式！！ 可參考 https://youtu.be/Qde65PVTR4I"


		case (t.Sub(lastNewbie) > cdNewbie && strings.Contains(msg,"是萌新")):
			lastNewbie = t
			replyMsg = "萌新是在說我嗎喵~ (探頭"
		case (t.Sub(lastBank) > cdBank && strings.Contains(msg,"池田銀行")):
			lastBank = t
			replyMsg = "點數太多嗎？歡迎存點數進來悠喵OwO"
		case (t.Sub(lastTeachMe) > cdTeachMe && teachMe(msg)):
			lastTeachMe = t
			replyMsg = appendSkyInfo(replyMsg)
//			replyMsg = appendStarflyxInfo(replyMsg+"\n\n")
			replyMsg = appendNTUSlidesInfo(replyMsg+"\n\n")
//			replyMsg = appendTaiwancoInfo(replyMsg+"\n\n")
//			replyMsg += "\n\nhttp://jmj.tw\n左上角還有些教學可以看 請多加利用喔喵~"
		case (strings.Contains(msg,"!天鳳") ):
			replyMsg = "天鳳位置在 https://tenhou.net/0/\n"+
						"各種說明可以在 http://tenhou.net/man/ 找到"
		case (strings.Contains(msg,"!雀姬")):
			replyMsg = "雀姬可以從 https://goo.gl/dQJFSm 下載\n"+
						"是手機上的遊戲喔~"
		case (strings.Contains(msg,"!雀魂") ):
			replyMsg = "雀魂位置在 https://www.majsoul.com/\n" +
						"用瀏覽器遊玩, 是個玩起來很舒服的遊戲~"
		case (strings.Contains(msg,"!戰績") ):
			replyMsg = "戰績網位置在 https://nodocchi.moe/tenhoulog\n"+
						"可以查到過往的戰績, 上方有過濾選項可以看特定時間或桌種"
		case (strings.Contains(msg,"!討論") ):
			replyMsg = "請大家跟著這樣做喔喵~\n\n" +

					"○ 等其他討論結束再開新議題\n"+
					"○ 用客觀的線索, 完整描述自己的想法、看法\n"+
					"○ 想想別人看到自己的發言, 會有什麼感覺\n"+
					"○ 講完自己的觀點後, 也看看別人怎麼說\n"+
					"○ 以善意的角度解讀其他人的發言\n"+
					"○ 覺得別人能做得更好時, 請提供你的看法做法\n"+
					"○ 不要害怕被否定或不被認同, 學得到東西就值得了\n\n"+

					"✕ 避免情緒性發言, 生氣不能解決問題\n" +
					"✕ 避免容易讓人誤解的玩笑\n"+
					"✕ 避免讓人難受或生氣說話方式\n"+
					"✕ 避免放大絕直接否定他人或中斷討論, 好好講就好了\n"

		case (strings.Contains(msg,"!何切") ):
			replyMsg ="「何切」是對一個既定場況探討該如何選擇/行動的討論方式\n\n" +
			"適合討論的場況需要：\n"+
			"✔是摸進手牌後、或可以選擇是否鳴牌的情境\n"+
			"✔遮擋他家手牌\n"+
			"✔認得出捨牌是否為摸打\n"+
			"✔知道鳴牌的時間點\n\n"+
			"討論時一般用0代表赤牌 m/p/s代表萬/筒/索\n"+
			"先講自己的選擇, 再用對場況的解讀和一些基於客觀線索的判斷作補充說明\n"+
			"儘量講得精確一點, 日麻變數太多, 籠統帶過別人可能無法瞭解你的意思by小劉\n\n"+
			"如果要問參考的何切分析 請用 [何切 1112345678999m1z] 這種方式詢問\n" +
			"後方接14個數字代表14張牌 mpsz代表花色 z是字牌喔~\n"+
			"※詳細說明請用「何切使用說明」詢問"
		case (strings.Contains(msg,"何切")):
			msg = strings.Replace(msg, " ", "", -1)
			words := strings.Fields(msg)
			re := regexp.MustCompile("(([0-9]+[MmPpSs])|([1-7]+[Zz]))+")
			result := "";
			reply := "";
			status := 0;
			for i := 0; i < len(words); i++ {
				result = re.FindString(words[i]);
				result = strings.ToLower(result);
				numAmount := 0;

				for j := 0 ; j < len(result) ; j++ {
					_, err := strconv.ParseFloat(string(result[j]), 64)
					if(err == nil) {numAmount++}
				}
				if numAmount > 6 && numAmount < 15{
					if numAmount%3 == 0 {
						status = -1
						reply = "張數不對 這樣不能拿去問天鳳姬呢喵~ \n(需要說明嗎? 請用「何切使用說明」詢問)"
					} else {
						countAry := make([]int, 35)
						for j := 0; j <= 34; j++ {
							countAry[j] = 0
						}

						pointer := 0
						countAmtAvailable := true;
						for j := len(result) - 1; j >= 0; j-- {
							if result[j] == 'z' {
								pointer = 27
							} else if result[j] == 's' {
								pointer = 18
							} else if result[j] == 'p' {
								pointer = 9
							} else if result[j] == 'm' {
								pointer = 0
							} else {
								num := int(result[j] - '0')
								if num == 0 {
									num = 5
								}
								countAry[num+pointer]++
								if countAry[num+pointer] > 4 {
									status = -1
									reply = "每種牌只能有四張喔！\n(需要說明嗎? 請用「何切使用說明」詢問)"
									countAmtAvailable = false;
									break;
								}
							}
						}
						if(!countAmtAvailable ) {continue;}

						if(numAmount%3 == 1)  {
							filled := false;
							for j := 1 ; j <= 7 ; j++ {
								if countAry[27+j] == 0 {
									result += strconv.Itoa(j) + "z";
									filled = true;
									break;
								}
							}
							if filled == false {
								kind := "mps";
								for j := 1 ; j <= 27 ; j++ {
									for k:= -2 ; k <= 2 ; k++ {
										if j+k < 1 {continue;}
										if k < 0 && (j + k-1) % 9 > (j-1) % 9 {continue;}
										if k > 0 && (j + k -1) % 9 < (j-1) % 9 {continue;}

										if countAry[j+k] > 0 {break;}
										if k == 2 {
											result +=  strconv.Itoa(((j-1)%9)+1) + string(kind[j/9]);
											filled = true;
											break;
										}
									}
									if filled == true {break; }
								}
							}
						}
						status = 1
						reply = "天鳳姬是這樣說的呢喵~\n" + "http://tenhou.net/2/?q=" + result + "\n";
						vm := otto.New()
						vm.Set("queryString", result)
						vm.Set("shantinInfo","")
						vm.Set("result","")
vm.Run(`

var u=function(){function b(a){var b=a&7,c=0,d=0;1==b||4==b?c=d=1:2==b&&(c=d=2);a>>=3;b=(a&7)-c;if(0>b)return!1;c=d;d=0;1==b||4==b?(c+=1,d+=1):2==b&&(c+=2,d+=2);a>>=3;b=(a&7)-c;if(0>b)return!1;c=d;d=0;1==b||4==b?(c+=1,d+=1):2==b&&(c+=2,d+=2);a>>=3;b=(a&7)-c;if(0>b)return!1;c=d;d=0;1==b||4==b?(c+=1,d+=1):2==b&&(c+=2,d+=2);a>>=3;b=(a&7)-c;if(0>b)return!1;c=d;d=0;1==b||4==b?(c+=1,d+=1):2==b&&(c+=2,d+=2);a>>=3;b=(a&7)-c;if(0>b)return!1;c=d;d=0;1==b||4==b?(c+=1,d+=1):2==b&&(c+=2,d+=2);a>>=3;b=(a&7)-c;if(0>b)return!1;c=d;d=0;1==b||4==b?(c+=1,d+=1):2==b&&(c+=2,d+=2);a>>=3;b=(a&7)-c;if(0!=b&&3!=b)return!1;b=(a>>3&7)-d;return 0==b||3==b}function a(a,d){if(0==a){if(128<=(d&448)&&b(d-128)||65536<=(d&229376)&&b(d-65536)||33554432<=(d&117440512)&&b(d-33554432))return!0}else if(1==a){if(16<=(d&56)&&b(d-16)||8192<=(d&28672)&&b(d-8192)||4194304<=(d&14680064)&&b(d-4194304))return!0}else if(2==a&&(2<=(d&7)&&b(d-2)||1024<=(d&3584)&&b(d-1024)||524288<=(d&1835008)&&b(d-524288)))return!0;return!1}function g(a,b){return a[b+0]<<0|a[b+1]<<3|a[b+2]<<6|a[b+3]<<9|a[b+4]<<12|a[b+5]<<15|a[b+6]<<18|a[b+7]<<21|a[b+8]<<24}function d(c){var d=1<<c[27]|1<<c[28]|1<<c[29]|1<<c[30]|1<<c[31]|1<<c[32]|1<<c[33];if(16<=d)return!1;if(2==(d&3)&&2==c[0]*c[8]*c[9]*c[17]*c[18]*c[26]*c[27]*c[28]*c[29]*c[30]*c[31]*c[32]*c[33]||!(d&10)&&7==(2==c[0])+(2==c[1])+(2==c[2])+(2==c[3])+(2==c[4])+(2==c[5])+(2==c[6])+(2==c[7])+(2==c[8])+(2==c[9])+(2==c[10])+(2==c[11])+(2==c[12])+(2==c[13])+(2==c[14])+(2==c[15])+(2==c[16])+(2==c[17])+(2==c[18])+(2==c[19])+(2==c[20])+(2==c[21])+(2==c[22])+(2==c[23])+(2==c[24])+(2==c[25])+(2==c[26])+(2==c[27])+(2==c[28])+(2==c[29])+(2==c[30])+(2==c[31])+(2==c[32])+(2==c[33]))return!0;if(d&2)return!1;var r=c[0]+c[3]+c[6],m=c[1]+c[4]+c[7],n=c[9]+c[12]+c[15],e=c[10]+c[13]+c[16],q=c[18]+c[21]+c[24],k=c[19]+c[22]+c[25],p=(r+m+(c[2]+c[5]+c[8]))%3;if(1==p)return!1;var l=(n+e+(c[11]+c[14]+c[17]))%3;if(1==l)return!1;var t=(q+k+(c[20]+c[23]+c[26]))%3;if(1==t||1!=(2==p)+(2==l)+(2==t)+(2==c[27])+(2==c[28])+(2==c[29])+(2==c[30])+(2==c[31])+(2==c[32])+(2==c[33]))return!1;r=(1*r+2*m)%3;m=g(c,0);n=(1*n+2*e)%3;e=g(c,9);q=(1*q+2*k)%3;c=g(c,18);return d&4?!(p|r|l|n|t|q)&&b(m)&&b(e)&&b(c):2==p?!(l|n|t|q)&&b(e)&&b(c)&&a(r,m):2==l?!(t|q|p|r)&&b(c)&&b(m)&&a(n,e):2==t?!(p|r|l|n)&&b(m)&&b(e)&&a(q,c):!1}return function(a,b){if(34==b)return d(a)}}();function w(){this.h=[-1,-1,-1,-1,-1,-1,-1];this.c=[{b:-1,a:0},{b:-1,a:0},{b:-1,a:0},{b:-1,a:0}]}w.prototype={};function x(b,a,g,d){b=b.c;var c=b[0].a,f=[0,0,0],r=7<<24-3*a,m=2<<24-3*a,n=0;(d&r)>=m&&y(c,g,d-m,f)&&(f[0]&&(b[n].b=9*g+8-a,b[n].a=f[0],++n),f[1]&&(b[n].b=9*g+8-a,b[n].a=f[1],++n),f[2]&&(b[n].b=9*g+8-a,b[n].a=f[2],++n));r>>=9;m>>=9;(d&r)>=m&&y(c,g,d-m,f)&&(f[0]&&(b[n].b=9*g+5-a,b[n].a=f[0],++n),f[1]&&(b[n].b=9*g+5-a,b[n].a=f[1],++n),f[2]&&(b[n].b=9*g+5-a,b[n].a=f[2],++n));m>>=9;(d&r>>9)>=m&&y(c,g,d-m,f)&&(f[0]&&(b[n].b=9*g+2-a,b[n].a=f[0],++n),f[1]&&(b[n].b=9*g+2-a,b[n].a=f[1],++n),f[2]&&(b[n].b=9*g+2-a,b[n].a=f[2],++n));return 0!=n}function z(b,a,g){b=b.c;var d=[0,0,0];if(!y(b[0].a,a,g,d))return!1;a=0;d[0]&&(b[a].b=b[0].b,b[a].a=d[0],++a);d[1]&&(b[a].b=b[0].b,b[a].a=d[1],++a);d[2]&&(b[a].b=b[0].b,b[a].a=d[2],++a);return 0!=a}function y(b,a,g,d){var c=-1,f,r=g&7,m=0,n=0;for(f=0;7>f&&1755!=g;++f){switch(r){case 4:b<<=8,b|=7*a+f+1,m+=1,n+=1;case 3:(g>>3&7)>=3+m&&(g>>6&7)>=3+n?(c=f,m+=3,n+=3):(b<<=8,b|=21+9*a+f+1);break;case 2:b<<=16;b|=257*(7*a+f+1);m+=2;n+=2;break;case 1:b<<=8;b|=7*a+f+1;m+=1;n+=1;break;case 0:break;default:return 0}g>>=3;r=(g&7)-m;m=n;n=0}if(7>f)return d[0]=16843009*(21+9*a+f+1)+66051,d[1]=65793*(7*a+f+1+1)|21+9*a+f+0+1<<24,d[2]=65793*(7*a+f+0+1)|21+9*a+f+3+1<<24,3;if(3==r)b=b<<8|9*a+29;else if(r)return 0;r=(g>>3&7)-m;if(3==r)b=b<<8|9*a+30;else if(r)return 0;if(-1!=c)return b<<=24,d[0]=b|65793*(21+9*a+c+1)+258,d[1]=b|65793*(7*a+c+1),d[2]=0,2;d[0]=b;d[1]=d[2]=0;return 1}function A(b,a,g,d){var c=7<<24-3*a,f=2<<24-3*a;if((d&c)>=f&&B(b,g,d-f))return b.c[0].b=9*g+8-a,!0;c>>=9;f>>=9;if((d&c)>=f&&B(b,g,d-f))return b.c[0].b=9*g+5-a,!0;f>>=9;return(d&c>>9)>=f&&B(b,g,d-f)?(b.c[0].b=9*g+2-a,!0):!1}function B(b,a,g){var d=b.c[0].a,c,f=g&7,r=0,m=0;for(c=0;7>c;++c){switch(f){case 4:d<<=16;d|=21+9*a+c+1<<8|7*a+c+1;r+=1;m+=1;break;case 3:d<<=8;d|=21+9*a+c+1;break;case 2:d<<=16;d|=257*(7*a+c+1);r+=2;m+=2;break;case 1:d<<=8;d|=7*a+c+1;r+=1;m+=1;break;case 0:break;default:return!1}g>>=3;f=(g&7)-r;r=m;m=0}if(3==f)d=d<<8|9*a+29;else if(f)return!1;f=(g>>3&7)-r;if(3==f)d=d<<8|9*a+30;else if(f)return!1;b.c[0].a=d;return!0}function C(b,a){var g,d=b.c,c=1<<a[27]|1<<a[28]|1<<a[29]|1<<a[30]|1<<a[31]|1<<a[32]|1<<a[33];if(16<=c)return!1;if(2==(c&3)&&2==a[0]*a[8]*a[9]*a[17]*a[18]*a[26]*a[27]*a[28]*a[29]*a[30]*a[31]*a[32]*a[33]){var f,c=[0,8,9,17,18,26,27,28,29,30,31,32,33];for(f=0;13>f&&2!=a[c[f]];++f);d[0].b=c[f];d[0].a=4294967295;return!0}if(c&2)return!1;f=!1;if(!(c&10)&&7==(2==a[0])+(2==a[1])+(2==a[2])+(2==a[3])+(2==a[4])+(2==a[5])+(2==a[6])+(2==a[7])+(2==a[8])+(2==a[9])+(2==a[10])+(2==a[11])+(2==a[12])+(2==a[13])+(2==a[14])+(2==a[15])+(2==a[16])+(2==a[17])+(2==a[18])+(2==a[19])+(2==a[20])+(2==a[21])+(2==a[22])+(2==a[23])+(2==a[24])+(2==a[25])+(2==a[26])+(2==a[27])+(2==a[28])+(2==a[29])+(2==a[30])+(2==a[31])+(2==a[32])+(2==a[33])){d[3].a=4294967295;for(f=g=0;34>f;++f)2==a[f]&&(b.h[g]=f,g+=1);f=!0}var r=a[0]+a[3]+a[6],m=a[1]+a[4]+a[7],n=a[2]+a[5]+a[8],e=a[9]+a[12]+a[15],q=a[10]+a[13]+a[16],k=a[11]+a[14]+a[17],p=a[18]+a[21]+a[24],l=a[19]+a[22]+a[25],t=a[20]+a[23]+a[26];g=(r+m+n)%3;if(1==g)return f;var v=(e+q+k)%3;if(1==v)return f;var h=(p+l+t)%3;if(1==h||1!=(2==g)+(2==v)+(2==h)+(2==a[27])+(2==a[28])+(2==a[29])+(2==a[30])+(2==a[31])+(2==a[32])+(2==a[33]))return f;c&8&&(3==a[27]&&(d[0].a<<=8,d[0].a|=49),3==a[28]&&(d[0].a<<=8,d[0].a|=50),3==a[29]&&(d[0].a<<=8,d[0].a|=51),3==a[30]&&(d[0].a<<=8,d[0].a|=52),3==a[31]&&(d[0].a<<=8,d[0].a|=53),3==a[32]&&(d[0].a<<=8,d[0].a|=54),3==a[33]&&(d[0].a<<=8,d[0].a|=55));n=r+m+n;r=(1*r+2*m)%3;m=D(a,0);k=e+q+k;e=(1*e+2*q)%3;q=D(a,9);t=p+l+t;p=(1*p+2*l)%3;l=D(a,18);if(c&4){if(g|r|v|e|h|p)return f;2==a[27]?d[0].b=27:2==a[28]?d[0].b=28:2==a[29]?d[0].b=29:2==a[30]?d[0].b=30:2==a[31]?d[0].b=31:2==a[32]?d[0].b=32:2==a[33]&&(d[0].b=33);if(9<=n){if(B(b,1,q)&&B(b,2,l)&&z(b,0,m))return!0}else if(9<=k){if(B(b,2,l)&&B(b,0,m)&&z(b,1,q))return!0}else if(9<=t){if(B(b,0,m)&&B(b,1,q)&&z(b,2,l))return!0}else if(B(b,0,m)&&B(b,1,q)&&B(b,2,l))return!0}else if(2==g){if(v|e|h|p)return f;if(8<=n){if(B(b,1,q)&&B(b,2,l)&&x(b,r,0,m))return!0}else if(9<=k){if(B(b,2,l)&&A(b,r,0,m)&&z(b,1,q))return!0}else if(9<=t){if(A(b,r,0,m)&&B(b,1,q)&&z(b,2,l))return!0}else if(B(b,1,q)&&B(b,2,l)&&A(b,r,0,m))return!0}else if(2==v){if(h|p|g|r)return f;if(8<=k){if(B(b,2,l)&&B(b,0,m)&&x(b,e,1,q))return!0}else if(9<=t){if(B(b,0,m)&&A(b,e,1,q)&&z(b,2,l))return!0}else if(9<=n){if(A(b,e,1,q)&&B(b,2,l)&&z(b,0,m))return!0}else if(B(b,2,l)&&B(b,0,m)&&A(b,e,1,q))return!0}else if(2==h){if(g|r|v|e)return f;if(8<=t){if(B(b,0,m)&&B(b,1,q)&&x(b,p,2,l))return!0}else if(9<=n){if(B(b,1,q)&&A(b,p,2,l)&&z(b,0,m))return!0}else if(9<=k){if(A(b,p,2,l)&&B(b,0,m)&&z(b,1,q))return!0}else if(B(b,0,m)&&B(b,1,q)&&A(b,p,2,l))return!0}d[0].a=0;return f}function D(b,a){return b[a+0]<<0|b[a+1]<<3|b[a+2]<<6|b[a+3]<<9|b[a+4]<<12|b[a+5]<<15|b[a+6]<<18|b[a+7]<<21|b[a+8]<<24};var E=function(){function b(a){e[a]-=2;++p}function a(a){e[a]+=2;--p}function g(a){--e[a];--e[a+1];--e[a+2];++q}function d(a){++e[a];++e[a+1];++e[a+2];--q}function c(a){--e[a];--e[a+1];++k}function f(a){++e[a];++e[a+1];--k}function r(a){--e[a];--e[a+2];++k}function m(a){++e[a];++e[a+2];--k}var n=0,e,q=0,k=0,p=0,l=0,t=0,v=0;return{g:8,v:function(){var a=8-2*q-k-p,b=q+k;p?b+=p-1:t&&v&&(t|v)==t&&++a;4<b&&(a+=b-4);-1!=a&&a<l&&(a=l);a<this.g&&(this.g=a)},l:function(a,b){e=[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0];v=t=l=p=k=q=0;this.g=8;if(136==b)for(b=0;136>b;++b)a[b]&&++e[b>>2];else if(34==b)for(b=0;34>b;++b)e[b]=a[b];else for(--b;0<=b;--b)++e[a[b]>>2]},j:function(){return e[0]+e[1]+e[2]+e[3]+e[4]+e[5]+e[6]+e[7]+e[8]+e[9]+e[10]+e[11]+e[12]+e[13]+e[14]+e[15]+e[16]+e[17]+e[18]+e[19]+e[20]+e[21]+e[22]+e[23]+e[24]+e[25]+e[26]+e[27]+e[28]+e[29]+e[30]+e[31]+e[32]+e[33]},o:function(){var a=(2<=e[0])+(2<=e[8])+(2<=e[9])+(2<=e[17])+(2<=e[18])+(2<=e[26])+(2<=e[27])+(2<=e[28])+(2<=e[29])+(2<=e[30])+(2<=e[31])+(2<=e[32])+(2<=e[33]),b=(0!=e[0])+(0!=e[8])+(0!=e[9])+(0!=e[17])+(0!=e[18])+(0!=e[26])+(0!=e[27])+(0!=e[28])+(0!=e[29])+(0!=e[30])+(0!=e[31])+(0!=e[32])+(0!=e[33]),d=b+(0!=e[1])+(0!=e[2])+(0!=e[3])+(0!=e[4])+(0!=e[5])+(0!=e[6])+(0!=e[7])+(0!=e[10])+(0!=e[11])+(0!=e[12])+(0!=e[13])+(0!=e[14])+(0!=e[15])+(0!=e[16])+(0!=e[19])+(0!=e[20])+(0!=e[21])+(0!=e[22])+(0!=e[23])+(0!=e[24])+(0!=e[25]),c=this.g,d=6-(a+(2<=e[1])+(2<=e[2])+(2<=e[3])+(2<=e[4])+(2<=e[5])+(2<=e[6])+(2<=e[7])+(2<=e[10])+(2<=e[11])+(2<=e[12])+(2<=e[13])+(2<=e[14])+(2<=e[15])+(2<=e[16])+(2<=e[19])+(2<=e[20])+(2<=e[21])+(2<=e[22])+(2<=e[23])+(2<=e[24])+(2<=e[25]))+(7>d?7-d:0);d<c&&(c=d);d=13-b-(a?1:0);d<c&&(c=d);return c},m:function(a){var b=0,d=0,c;for(c=27;34>c;++c)switch(e[c]){case 4:++q;b|=1<<c-27;d|=1<<c-27;++l;break;case 3:++q;break;case 2:++p;break;case 1:d|=1<<c-27}l&&2==a%3&&--l;d&&(v|=134217728,(b|d)==b&&(t|=134217728))},w:function(a){var b=0,d=0,c;for(c=27;34>c;++c)switch(e[c]){case 4:++q;b|=1<<c-18;d|=1<<c-18;++l;break;case 3:++q;break;case 2:++p;break;case 1:d|=1<<c-18}for(c=0;9>c;c+=8)switch(e[c]){case 4:++q;b|=1<<c;d|=1<<c;++l;break;case 3:++q;break;case 2:++p;break;case 1:d|=1<<c}l&&2==a%3&&--l;d&&(v|=134217728,(b|d)==b&&(t|=134217728))},s:function(a){t|=(4==e[0])<<0|(4==e[1])<<1|(4==e[2])<<2|(4==e[3])<<3|(4==e[4])<<4|(4==e[5])<<5|(4==e[6])<<6|(4==e[7])<<7|(4==e[8])<<8|(4==e[9])<<9|(4==e[10])<<10|(4==e[11])<<11|(4==e[12])<<12|(4==e[13])<<13|(4==e[14])<<14|(4==e[15])<<15|(4==e[16])<<16|(4==e[17])<<17|(4==e[18])<<18|(4==e[19])<<19|(4==e[20])<<20|(4==e[21])<<21|(4==e[22])<<22|(4==e[23])<<23|(4==e[24])<<24|(4==e[25])<<25|(4==e[26])<<26;q+=a;this.u(0)},u:function(h){var k=arguments.callee;++n;if(-1!=this.g){for(;27>h&&!e[h];++h);if(27==h)return this.v();var l=h;8<l&&(l-=9);8<l&&(l-=9);switch(e[h]){case 4:e[h]-=3;++q;7>l&&e[h+2]&&(e[h+1]&&(g(h),k.call(this,h+1),d(h)),r(h),k.call(this,h+1),m(h));8>l&&e[h+1]&&(c(h),k.call(this,h+1),f(h));var p=h;--e[p];v|=1<<p;k.call(this,h+1);p=h;++e[p];v&=~(1<<p);e[h]+=3;--q;b(h);7>l&&e[h+2]&&(e[h+1]&&(g(h),k.call(this,h),d(h)),r(h),k.call(this,h+1),m(h));8>l&&e[h+1]&&(c(h),k.call(this,h+1),f(h));a(h);break;case 3:e[h]-=3;++q;k.call(this,h+1);e[h]+=3;--q;b(h);7>l&&e[h+1]&&e[h+2]?(g(h),k.call(this,h+1),d(h)):(7>l&&e[h+2]&&(r(h),k.call(this,h+1),m(h)),8>l&&e[h+1]&&(c(h),k.call(this,h+1),f(h)));a(h);7>l&&2<=e[h+2]&&2<=e[h+1]&&(g(h),g(h),k.call(this,h),d(h),d(h));break;case 2:b(h);k.call(this,h+1);a(h);7>l&&e[h+2]&&e[h+1]&&(g(h),k.call(this,h),d(h));break;case 1:6>l&&1==e[h+1]&&e[h+2]&&4!=e[h+3]?(g(h),k.call(this,h+2),d(h)):(p=h,--e[p],v|=1<<p,k.call(this,h+1),p=h,++e[p],v&=~(1<<p),7>l&&e[h+2]&&(e[h+1]&&(g(h),k.call(this,h+1),d(h)),r(h),k.call(this,h+1),m(h)),8>l&&e[h+1]&&(c(h),k.call(this,h+1),f(h)))}}}}}();function F(b,a){E.l(b,34);var g=E.j();if(14<g)return-2;!a&&13<=g&&(E.g=E.o(g));E.m(g);E.s(Math.floor((14-g)/3));return E.g}function G(b,a){E.l(b,a);var g=E.j();if(!(14<g)){var d=[E.g,E.g];13<=g&&(d[0]=E.o(g));E.m(g);E.s(Math.floor((14-g)/3));d[1]=E.g;d[1]<d[0]&&(d[0]=d[1]);return d}};function H(b){var a=b>>2;return(27>a&&16==b%36?"0":a%9+1)+"mpsz".substr(a/9,1)}function J(b){return b.replace(/(m|p|s|z)/g,"$&:").replace(/(m|p|s|z)([^:])/g,"$2").replace(/:/g,"")}function aa(b){b=b.replace(/(\d)m/g,"0$1").replace(/(\d)p/g,"1$1").replace(/(\d)s/g,"2$1").replace(/(\d)z/g,"3$1");var a,g=Array(136);for(a=0;a<b.length;a+=2){var d=b.substr(a,2),c;d%10?(c=4*(9*Math.floor(d/10)+(d%10-1)),c=g[c+3]?g[c+2]?g[c+1]?c:c+1:c+2:c+3):c=4*(9*d/10+4)+0;g[c]&&document.write("err n="+d+" k="+c+"");g[c]=1}return g};function ba(b){var a=parseInt(b.substr(0,1));return(a?a-1:4)+9*"mpsz".indexOf(b.substr(1,1))}function K(b){var a,g=[];for(a=0;34>a;++a)4<=b[a]||(b[a]++,u(b,34)&&g.push(a),b[a]--);return g}function ca(b){var a,g=[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0];for(a=0;136>a;++a)b[a]&&++g[a>>2];return g}function da(b,a,g,d,c,f){}function M(b){return-1==b?"\u548c\u4e86":0==b?"\u8074\u724c":b+"\u5411\u8074"}function N(b,a){return a&&b[0]!=b[1]?"\u6a19\u6e96\u5f62"+M(b[0])+" / \u4e00\u822c\u5f62"+M(b[1]):M(b[0])}function ea(b){function a(a){a&=127;return 21>a?(a=9*Math.floor(a/7)+a%7,L(H(4*a+1))+L(H(4*a+5))+L(H(4*a+9))):55>a?(a=L(H(4*(a-21)+1)),a+a+a):89>a?(a=L(H(4*(a-55)+1)),a+a+a+a):""}function g(a){a=L(H(4*a+1));return a+a}var d=new w;if(C(d,b)){var c="";for(b=0;4>b;++b)if(d.c[b].a){if(0==b&&4294967295==d.c[0].a){var c=c+"\u56fd\u58eb\u5f62\u548c\u4e86 ",c=c+(g(d.c[b].b)+" "),f,r=[0,8,9,17,18,26,27,28,29,30,31,32,33];for(f=0;13>f;++f)d.c[b].b!=r[f]&&(c+=L(H(4*r[f]+1)))}else 3==b&&4294967295==d.c[3].a?(c+="\u4e03\u5bfe\u5f62\u548c\u4e86 ",c+=g(d.h[0])+" "+g(d.h[1])+" "+g(d.h[2])+" "+g(d.h[3])+" "+g(d.h[4])+" "+g(d.h[5])+" "+g(d.h[6])):(f=[(d.c[b].a>>0&255)-1,(d.c[b].a>>8&255)-1,(d.c[b].a>>16&255)-1,(d.c[b].a>>24&255)-1],c+="\u4e00\u822c\u5f62\u548c\u4e86 ",c+=g(d.c[b].b)+" "+a(f[3])+" "+a(f[2])+" "+a(f[1])+" "+a(f[0]));c+=""}return c}}function fa(){var O=queryString,O=O.replace(/^([^=]+)=(.+)/,"$2"),ga="q";function b(a,b){var c,d=0;for(c=0;c<a.length;++c)d+=4-b[a[c]];return d}var a=ga,g=O,d;d="";switch(a.substr(0,1)){case "q":d+="\u6a19\u6e96\u5f62(\u4e03\u5bfe\u56fd\u58eb\u3092\u542b\u3080)\u306e\u8a08\u7b97\u7d50\u679c / "+a.substr(1)+"="+g+"\u4e00\u822c\u5f62";break;case "p":d+="\u4e00\u822c\u5f62(\u4e03\u5bfe\u56fd\u58eb\u3092\u542b\u307e\u306a\u3044)\u306e\u8a08\u7b97\u7d50\u679c / "+a.substr(1)+"="+g+"\u6a19\u6e96\u5f62"}for(var c="d"==a.substr(1,1),a=a.substr(0,1),g=g.replace(/(\d)(\d{0,8})(\d{0,8})(\d{0,8})(\d{0,8})(\d{0,8})(\d{0,8})(\d{8})(m|p|s|z)/g,"$1$9$2$9$3$9$4$9$5$9$6$9$7$9$8$9").replace(/(\d?)(\d?)(\d?)(\d?)(\d?)(\d?)(\d)(\d)(m|p|s|z)/g,"$1$9$2$9$3$9$4$9$5$9$6$9$7$9$8$9").replace(/(m|p|s|z)(m|p|s|z)+/g,"$1").replace(/^[^\d]/,""),g=g.substr(0,28),f=aa(g),r=-1;r=Math.floor(136*Math.random()),f[r];);var m=Math.floor(g.length/2)%3;2==m||c||(f[r]=1,g+=H(r));var f=ca(f),n="",e=G(f,34),n=n+N(e,28==g.length);for(var i=0;i<n.length;i++)shantinInfo+=n[i];n=n+("("+Math.floor(g.length/2)+"\u679a)");-1==e[0]&&(n+=" / \u65b0\u3057\u3044\u624b\u724c\u3092\u4f5c\u6210");var n=n+"",q="q"==a?e[0]:e[1],k,p,l=Array(35);if(0==q&&1==m&&c)k=34,l[k]=K(f),l[k].length&&(l[k]={i:k,n:b(l[k],f),c:l[k]});else if(0>=q)for(k=0;34>k;++k)f[k]&&(f[k]--,l[k]=K(f),f[k]++,l[k].length&&(l[k]={i:k,n:b(l[k],f),c:l[k]}));else if(2==m||1==m&&!c)for(k=0;34>k;++k){if(f[k]){f[k]--;l[k]=[];for(p=0;34>p;++p)k==p||4<=f[p]||(f[p]++,F(f,"p"==a)==q-1&&l[k].push(p),f[p]--);f[k]++;l[k].length&&(l[k]={i:k,n:b(l[k],f),c:l[k]})}}else{k=34;l[k]=[];for(p=0;34>p;++p)4<=f[p]||(f[p]++,F(f,"p"==a)==q-1&&l[k].push(p),f[p]--);l[k].length&&(l[k]={i:k,n:b(l[k],f),c:l[k]})}var t=[];for(k=0;k<g.length;k+=2){p=g.substr(k,2);var v=ba(p),h=J(g.replace(p,"").replace(/(\d)(m|p|s|z)/g,"$2$1$1,").replace(/00/g,"50").split(",").sort().join("").replace(/(m|p|s|z)\d(\d)/g,"$2$1")),R=q+1,I=l[v];I&&I.n&&(R=-1==q?0:q,void 0==I.q&&t.push(I),I.q=h);2==m&&(h+=H(r));n+=(2==m||2!=m&&!c?da:L)(p,2==k%3&&k==g.length-2?"":"",a,h,v,R)}l[34]&&l[34].n&&(l[34].q=J(g),t.push(l[34]));t.sort(function(a,b){return b.n-a.n});g=""+(queryString+"\n");q=0>=q?"\u5f85\u3061":"\u6478";for(k=0;k<t.length;++k){v=t[k].i;34>v&&(g+="\u6253"+H(4*v+1)+" ");g+=q+"[";l=t[k].c;c=t[k].q;for(p=0;p<l.length;++p)r=H(4*l[p]+1),g+=H(4*l[p]+1);g+=" "+t[k].n+"\u679a]\n"};for(var i=0;i<g.length;i++)result+=g[i];}fa();
	`)
						if value, err := vm.Get("shantinInfo"); err == nil {
							a, _ := value.ToString();
							reply += a + "\n";
						}
						if value, err := vm.Get("result"); err == nil {
							a, _ := value.ToString();
							reply += a;
						}
/*	這邊沒用上...不是想要的內容
// https://www.daniweb.com/programming/computer-science/code/495192/get-the-content-of-a-web-page-golang
	url := "http://tenhou.net/2/?q=" + result;

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader([]byte(html))
    doc, _ := goquery.NewDocumentFromReader(r)
    text := doc.Find("textarea").Text()
    reply += text;
*/
						break;
					}
				}
			}
			if(status != 0) {replyMsg = reply}
			if(t.Sub(lastWhatCutHelp) > cdWhatCutHelp && _024_usage.MatchString(msg)) {
/*
				replyMsg = "手牌必須是數字接花色的組合 m萬p筒s索 z字\n" +
				"三色牌數字是0~9 其中0是赤\n" +
				"字牌的話數字只能用1-7\n" +
				"同一種牌最多只會有四張 不要自己刻不存在的牌嘿乖～\n\n"+
				"問何切要3n+2張牌 不然丟一張出去會相公喔！\n"+
				"問聽牌則是要3n+1張, 可以在後面多加一張無關牌"
*/
				replyMsg = "想問何切嗎？請給我8/11/14張牌吧！\n\n"+

				"m=萬 p=筒 s=索 z=字 0=赤\n"+
				"字牌只有1-7z 分別代表東南西北白發中\n\n"+

				"例：「何切 34067m46p6667s22z6p」是問這手牌\n"+
				"三四五六七四六六六六七南南 六\n"+
				"萬萬萬萬萬筒筒索索索索風風 筒\n"+
				"(五萬是紅的)\n\n"+

				"張數正確、且確實可能出現(一種牌最多四張)才能幫你解答喔~\n"+
				"如果少一張的話，喵會幫你補一張無關牌進去，記得感謝我喔喵~"
			}
		case (t.Sub(lastGiveUp) > cdGiveUp && strings.Contains(msg,"棄麻")) :
			lastGiveUp = t
			replyMsg = "棄麻"
		case (t.Sub(lastSlides) > cdSlides && askingNTUSlides(msg)) :
			lastSlides = t
			replyMsg = appendNTUSlidesInfo(replyMsg)
		case (!groupExcluded && _029_interact.MatchString(msg)):
			switch {
				case (_028_ikedang.MatchString(msg)):
				replyMsg = "欸？不可以亂來喔喵 > <"
				case (_027_ikedaok.MatchString(msg)):
				replyMsg = "(´,,•ω•,,)開心開心"
				default:
				replyMsg = "不是很懂這人是想做啥呢喵～？"
			}
		case(_030_askDragon.MatchString(msg)):
			replyMsg = appendSkyInfo(replyMsg)
		case(_031_askStarflyx.MatchString(msg)):
			replyMsg = appendStarflyxInfo(replyMsg)
		case(_032_askBingo.MatchString(msg)):
			replyMsg = appendNTUSlidesInfo(replyMsg)
		case(_033_askTaiwanco.MatchString(msg)):
			replyMsg = appendTaiwancoInfo(replyMsg)

/*其他遊戲用途*/
/*puyo tetris 上方是影片 下方是圖片*/
		case (strings.Contains(msg,"!氣泡教學")):
			replyMsg = "(魔法氣泡)目前支援名詞解說如下~\n"+
			"!連鎖尾 !GTR !KEY !31流 !折返"
		case (strings.Contains(msg,"!方塊教學")):
			replyMsg = "(俄羅斯方塊)目前支援名詞解說如下~\n"+
			"!DT砲 !PC !4w !TKI !b2b"
		case (strings.Contains(msg,"!dt")):
			replyMsg = "!dt"
		case (strings.Contains(msg,"!pc")):
			replyMsg = "!pc"
		case (strings.Contains(msg,"!c4w") || strings.Contains(msg,"!4w")):
			replyMsg = "!c4w"
		case (strings.Contains(msg,"!tki")):
			replyMsg = "!tki"
		case (strings.Contains(msg,"!連鎖尾")):
			replyMsg = "!連鎖尾"
		case (strings.Contains(msg,"!b2b")):
			replyMsg = "!b2b"
		case (strings.Contains(msg,"!gtr") || strings.Contains(msg,"!新gtr")):
			replyMsg = "!gtr"
		case (strings.Contains(msg,"!key")):
			replyMsg = "!key"
		case (strings.Contains(msg,"!階梯") || strings.Contains(msg,"!31流")):
			replyMsg = "!31流"
		case (strings.Contains(msg,"!折返")):
			replyMsg = "!折返"

/*黑白棋*/
		case (strings.Contains(msg,"!黑白棋教學")):
			replyMsg = "素材徵求中, 目前支援名詞解說如下, 感謝草草提供~\n"+
			"!偶數理論 !奇偶性 !餘裕手 !開放度 !機動性 !穩定子 !天王山 !逆轉奇偶 !不平衡邊 !平衡邊"
		case (strings.Contains(msg,"黑白棋")):
			switch {
				case (_025_wheretoplay.MatchString(msg)):
				replyMsg = "這裡可以下黑白棋喔 ~ http://wars.fm/reversi"
				case (_026_wheretoplayonmobile.MatchString(msg)):
				replyMsg = "手機上玩黑白棋嗎~? 試試這個吧 ~\n"+ "https://play.google.com/store/apps/details?id=fm.wars.reversi"
				default:
			}

		case (strings.Contains(msg,"!偶數理論") || strings.Contains(msg,"!奇偶性")):
			replyMsg = "!偶數理論"
		case (strings.Contains(msg,"!餘裕手")):
			replyMsg = "!餘裕手"
		case (strings.Contains(msg,"!開放度")):
			replyMsg = "!開放度"
		case (strings.Contains(msg,"!機動性")):
			replyMsg = "!機動性"
		case (strings.Contains(msg,"!穩定子")):
			replyMsg = "!穩定子"
		case (strings.Contains(msg,"!天王山")):
			replyMsg = "!天王山"
		case (strings.Contains(msg,"!逆轉奇偶")):
			replyMsg = "!逆轉奇偶"
		case (strings.Contains(msg,"!不平衡邊") || strings.Contains(msg,"!平衡邊")):
			replyMsg = "!平衡邊"
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
				/*把全型！轉成半型!*/
				message.Text = strings.Replace(message.Text, "！", "!", -1)
				/*把大寫轉小寫*/
				message.Text = strings.ToLower(message.Text)

				replyMsg := determineReply(message.Text, event.Source.Type == "group" && isSupportedGroup(event.Source.GroupID), event.Source.Type == "group" && isExcludedGroup(event.Source.GroupID))

				if replyMsg == "棄麻" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage("https://i.imgur.com/9kmdMYH.jpg", "https://i.imgur.com/9kmdMYH.jpg")).Do(); err != nil {
						log.Print(err)
					}
					return
				}

				if replyMsg == "欸？不可以亂來喔喵 > <" {
					reply1msg1pic(event.ReplyToken, "欸？不可以亂來喔喵 > <", "https://i.imgur.com/9Zy1CXe.jpg")
					return
				}
/*puyo tetris 名詞專區*/
				if replyMsg == "!dt" {
					reply1msg1video(event.ReplyToken,
					"[Tetris] DT砲：\n攻勢由兩個T轉組成，分別消2,3行，故被稱為Double Triple Cannon，兩個T轉合共能至少造成11行傷害。\n擁有搭建便利，爆發力高，後續接招方法廣泛的特點，是為目前最受歡迎的開局之一。",
					"https://i.imgur.com/F5lmXur.mp4",
					"https://i.imgur.com/U7CNgL5.png")
					return
				}
				if replyMsg == "!pc" {
					reply1msg1video(event.ReplyToken,
					"[Tetris] Perfect Clear / PC / 全消：\n當玩家把場上所有方塊都消除後會發動追加攻擊，造成額外10行傷害，為目前單發傷害最高的攻擊手段。\n",
					"https://i.imgur.com/I7bCJGz.mp4",
					"https://i.imgur.com/5PJrPfL.png")
					return
				}
				if replyMsg == "!c4w" {
					reply1msg1video(event.ReplyToken,
					"[Tetris] c4w / 4w / 4-wide：\n在中間留下一道4行寛的深坑，能承受大量攻擊，也能通過連擊的加成作出反擊。\n常出現被對手速攻壓制，或反擊秒殺對手的局面。\n也有將四行空格留於側邊的排法。留四統稱4w(4-wide)",
					"https://i.imgur.com/mCz46GS.mp4",
					"https://i.imgur.com/GNXzIoO.png")
					return
				}
				if replyMsg == "!tki" {
					reply1msg1video(event.ReplyToken,
					"[Tetris] TKI：\n以開局穩定作出大量T轉為目標的排法。\n有成功率高、易上手、好變化等優點",
					"https://i.imgur.com/CueyM1O.mp4",
					"https://i.imgur.com/7TwPyvH.png")
					return
				}
				if replyMsg == "!連鎖尾" {
					reply1msg1video(event.ReplyToken,
					"[Puyo]連鎖尾：\n在主要連鎖最後不好發展的部分多排一些氣泡，增加連鎖的長度和攻擊量。",
					"https://i.imgur.com/v4NgNqZ.mp4",
					"https://i.imgur.com/qGIT73p.png")
					return
				}
				if replyMsg == "!b2b" {
					reply1msg1pic(event.ReplyToken,
					"[Tetris]Back-To-Back / b2b：\n當玩家做出T轉或Tetris消除行數時會進入B2B狀態，\n在這狀態下消除T轉或Tetris時會追加1行攻擊",
					"https://i.imgur.com/pL0Rxm9.png")
					return
				}
				if replyMsg == "!gtr" {
					reply1msg1pic(event.ReplyToken,
					"[Puyo] GTR / 新GTR：\n魔法氣泡中很常見的折返(轉彎)排法，\n因為形狀較為扁平，上方比較容易延伸連鎖\n(GTR和新GTR的差別在轉彎處的寬度)",
					"https://i.imgur.com/byUTOEQ.png")
					return
				}
				if replyMsg == "!key" {
					reply1msg1pic(event.ReplyToken,
					"[Puyo] Key / 三明治 / 鑰匙排法：\n因為排列起來很像一層夾著一層的三明治，也很像鑰匙的形狀而得名。\n在排列上變化較為自由，可以與其他折返排法合併使用。",
					"https://i.imgur.com/RlELvNo.png")
					return
				}
				if replyMsg == "!31流" {
					reply1msg1pic(event.ReplyToken,"[Puyo] 31流 / 階梯排法：\n階梯排法的優點是直觀好上手，依照氣泡數量不同又分為3-1階梯、2-2階梯等，\n但在折返(轉彎)處，可能發生連鎖中誤消氣泡的情形，\n因此在往上發展時較為弱勢，也較難快速的排至第二層。",
					"https://i.imgur.com/irmkrlt.png")
					return
				}
				if replyMsg == "!折返" {
					reply1msg1pic(event.ReplyToken,
					"[Puyo] 折返 / 轉彎：\n魔法氣泡的排列空間只有六行，大連鎖需要排到第二層以上。\n我們稱往上排的地方為「折返」。",
					"https://i.imgur.com/Uxkmm9f.png")
					return
				}


/* 黑白棋名詞專區 */
				if replyMsg == "!偶數理論" {
					reply1msg1pic(event.ReplyToken,
					"偶數理論（奇偶性）：留給對手下子的每個區域內, 都留下偶數個空位的策略。\n(黑子下在黃星處後, 所有區域都剩偶數個空位)",
					"https://i.imgur.com/ifOGTjR.png")
					return
				}

				if replyMsg == "!餘裕手" {
					reply1msg1pic(event.ReplyToken,
					"餘裕手：不會為對手帶來新的落子選擇的好棋。\n(下例中黑子先下1再下2, 對手就會因為沒有其他選擇, 而被迫讓出角落)",
					"https://i.imgur.com/gezuvEV.png")
					return
				}

				if replyMsg == "!開放度" {
					reply1msg1pic(event.ReplyToken,					"開放度：被翻轉的棋子中，每一子周圍(八格)的空格數總和，開放度越小越好。\n(黑子下在標示處, 一共只翻轉一顆白子, 開放度總和為1)",
					"https://i.imgur.com/Vdhu1Si.png")
					return
				}

				if replyMsg == "!機動性" {
					reply1msg1pic(event.ReplyToken,
					"機動性(Mobility)：可以落子的地方。\n一般來說, 機動性越高越有利, 因為愈可能存在較佳的路線；\n相對地, 選擇少的時候, 就容易被對手逼死\n(打x的是黑棋可以選擇落子的地方, 這個例子的機動性是5)",
					"https://i.imgur.com/VNbe0Zj.png")
					return
				}

				if replyMsg == "!穩定子" {
					reply1msg1pic(event.ReplyToken,
					"穩定子：永遠不會被翻轉的棋。(如下圖中的黑子)",
					"https://i.imgur.com/W1GY2PZ.png")
					return
				}

				if replyMsg == "!天王山" {
					reply1msg1pic(event.ReplyToken,
					"天王山：對黑白兩方都有利的位置。\n(兩邊下在黃星處後, 都不會幫對手增加多少選擇)",
					"https://i.imgur.com/hv3SkiV.png")
					return
				}

				if replyMsg == "!逆轉奇偶" {
					reply1msg1pic(event.ReplyToken,					"逆轉奇偶：利用棄權，讓對手被迫先下入偶數區堿而陷入不利狀態的戰術。\n(白棋的選擇被限制了, 在這個區域內先落子的人不利)",
					"https://i.imgur.com/ADexsHc.png")
					return
				}

				if replyMsg == "!平衡邊" {
					reply1msg1pic(event.ReplyToken,					"不平衡邊：自己只有連續3~5顆棋的邊(紅框處)，這類棋型容易遭到對手攻擊。\n平衡邊：六子邊(藍框處)，這種棋型較為安全。",
					"https://i.imgur.com/lddPkl1.png")
					return
				}

/*
				if replyMsg == "測試" {
					if(isAdmin(event.Source.UserID) && event.Source.Type == "group") {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(event.Source.GroupID)).Do(); err != nil {
							log.Print(err)
						}
					}
					return
				}
*/
				if replyMsg != "" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMsg)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}
