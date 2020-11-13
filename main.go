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
			default: //天氣單項
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
			bot.Send(msg)
		}
		//bot喚醒
		if update.Message != nil {
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
