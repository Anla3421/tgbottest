package main

import (
	"fmt"
	"log"
	"server/db/sql"
	"server/movie"
	"server/webpage"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const JaredID int64 = 1366159494
const BotToken = "1307746334:AAEmCDB1-OdP25rMjOK30zFLjJA8psEUviI"

var Sugar string
var Arguments string
var Drinkid int = 0

//Inlinekeyboard Setting
var InitialKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("天氣", "天氣查詢"),
		tgbotapi.NewInlineKeyboardButtonData("電影", "豆瓣網電影精選"),
		tgbotapi.NewInlineKeyboardButtonURL("Test Page", "http://127.0.0.1:3388/test/1"),
	),
)
var WeatherKB = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("台北", "1"),
		tgbotapi.NewInlineKeyboardButtonData("台中", "2"),
		tgbotapi.NewInlineKeyboardButtonData("台南", "3"),
		tgbotapi.NewInlineKeyboardButtonData("高雄", "4"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("我全都要", "5"),
	),
)
var MovieKB tgbotapi.InlineKeyboardMarkup

func init() {
	var row = []tgbotapi.InlineKeyboardButton{}
	var total = [][]tgbotapi.InlineKeyboardButton{}
	for i := 0; i < 11; i++ {
		button := tgbotapi.NewInlineKeyboardButtonData("第"+strconv.Itoa(i)+"頁", strconv.Itoa(i))
		row = append(row, button)
		if i%3 == 0 {
			total = append(total, tgbotapi.NewInlineKeyboardRow(row...))
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	MovieKB = tgbotapi.NewInlineKeyboardMarkup(total...)
}

var SugarKB = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("全糖", "全糖"),
		tgbotapi.NewInlineKeyboardButtonData("半糖", "半糖"),
		tgbotapi.NewInlineKeyboardButtonData("微糖", "微糖"),
		tgbotapi.NewInlineKeyboardButtonData("無糖", "無糖"),
	),
)
var IceKB = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("全冰", "全冰"),
		tgbotapi.NewInlineKeyboardButtonData("半冰", "半冰"),
		tgbotapi.NewInlineKeyboardButtonData("微冰", "微冰"),
		tgbotapi.NewInlineKeyboardButtonData("去冰", "去冰"),
	),
)

//---------------------------------------------------------------------------
func main() {
	go webpage.StartWebServer()
	go movie.Moviespider()
	fmt.Println("bot initial")
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	fmt.Print("Success connected, bot online.")
	for update := range updates {
		if update.CallbackQuery != nil {

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			switch update.CallbackQuery.Data {
			//第二層面板喚醒
			case "天氣查詢":
				msg.ReplyMarkup = WeatherKB
			case "豆瓣網電影精選":
				msg.ReplyMarkup = MovieKB
			//天氣選項回應
			case "5": //我全都要
				msg.Text = ""
				for ID := 1; ID < 5; ID++ {
					result := sql.Weathersql(strconv.Itoa(ID)).Text + "\n"
					msg.Text += result
				}
			default: //天氣單項1~4
				msg.Text = ""
				result := sql.Weathersql(update.CallbackQuery.Data).Text
				if result != "" {
					msg.Text += result
				}

			}
			//電影選項回應
			if update.CallbackQuery.Message.Text == "豆瓣網電影精選" {
				ID, _ := strconv.Atoi(update.CallbackQuery.Data)
				msg.Text = ""
				for Rank := 25*ID + 1; Rank < 25*ID+26; Rank++ {
					result := sql.Moviesqlget(Rank).Idre + "  " + sql.Moviesqlget(Rank).Moviename + "\n"
					msg.Text += result
				}
				msg.Text = "電影推薦TOP" + strconv.Itoa(1+25*(ID)) + "~" + strconv.Itoa(25*(ID+1)) + ":\n" +
					msg.Text + "\n請使用https://movie.douban.com/subject/\n+上述電影的數字來進入該電影之介紹"
			}
			if update.CallbackQuery.Message.Text == "請選擇甜度" {
				msg.Text = update.CallbackQuery.Data
				Sugar = update.CallbackQuery.Data
				bot.Send(msg)
				msg.Text = "請選擇冰量"
				msg.ReplyMarkup = IceKB
			}
			if update.CallbackQuery.Message.Text == "請選擇冰量" {
				msg.Text = update.CallbackQuery.Data
				Ice := msg.Text
				bot.Send(msg)
				Drinkid = Drinkid + 1
				Who := update.CallbackQuery.Message.Chat.FirstName

				msg.Text = Who + "點了: " + Arguments + " " + Sugar + update.CallbackQuery.Data
				sql.Drinksql(Drinkid, Who, Arguments, Sugar, Ice)
			}
			bot.Send(msg)
		}
		//bot喚醒
		if update.Message != nil {
			if update.Message.IsCommand() {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				switch update.Message.Command() {
				case "help":
					msg.Text = "type /drink 飲料名. or /drink-飲料名. to order drink.\n/total, /clear."
				case "drink":
					Arguments = update.Message.CommandArguments()
					if Arguments == "" {
						msg.Text = "type /drink 飲料名. or /drink-飲料名."
					} else {
						msg.Text = "飲料: " + update.Message.CommandArguments() + "已點餐"
						bot.Send(msg)
						msg.Text = "請選擇甜度"
						msg.ReplyMarkup = SugarKB
					}

				case "total":
					msg.Text = ""
					for ID := 1; ID < Drinkid+1; ID++ {
						result := strconv.Itoa(sql.Drinksqlget(ID).ID) + "." + sql.Drinksqlget(ID).Who + "\t" +
							sql.Drinksqlget(ID).Drink + " " + sql.Drinksqlget(ID).Sugar + sql.Drinksqlget(ID).Ice + "\n"
						msg.Text += result
					}
				case "clear":
					if update.Message.Chat.ID == JaredID {
						sql.Drinksqltruncate()
						Drinkid = 0
						msg.Text = "table clear complete"
					} else {
						msg.Text = "您沒有權限執行這個指令"
					}

				default:
					msg.Text = "type /help"
				}
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

				switch update.Message.Text {
				case "hi":
					msg.ReplyMarkup = InitialKeyboard
					msg.Text = "您今天想要點什麼呢?"
				default:
					//對特定人士回復特定的問候語
					switch update.Message.Chat.ID {
					case JaredID:
						msg.Text = "主人您好"
					default:
						msg.Text = "指令錯誤"
					}
					bot.Send(tgbotapi.NewMessage(JaredID, "有人對你說，他是 "+strconv.Quote(update.Message.Chat.FirstName)+strconv.Quote(update.Message.Chat.LastName)+" 他說： "+update.Message.Text+", ID是 "+strconv.FormatInt(update.Message.Chat.ID, 10)))
				}
				bot.Send(msg)
			}
		}
	}
}
