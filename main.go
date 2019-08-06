package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)
//токен телеграм бота
var token string = ""
//id чата куда бот будет отправлять сообщения, необходимо предварительно его туда добавить
var chatId int64  = -0000000000
var startStr string = "Как разместить свое объявление на @channel_name \n\nВсе просто:\n- Пишите /start\n- Выбираете категорию \"Сдам\" или \"Сниму\"\n- Введите текст объявления. Не забудьте указать адрес, станцию метро и стоимость."
var helpStr string = "/start начать диалог с @channel_name_bot\n/help - список всех команд"

var mainMenu1 = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🏠 Сдам"),
		tgbotapi.NewKeyboardButton("🚶🏻‍♂️ Сниму"),
	),
)

var mainMenuBack = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("❎ Назад"),
	),
)
var mainMenuBackAndSkip = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("❎ Назад"),
		tgbotapi.NewKeyboardButton("😎 Пропустить"),
	),
)
var mainMenuComplete = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("✅ Готово"),
	),
)
type probsummarySing struct {
	State int 
	Category string
	Text string
	Photo []string
}

var probsummarySingMap map[int]*probsummarySing

func init() {
	probsummarySingMap = make(map[int]*probsummarySing)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	//log.Printf("Авторизация %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if (update.Message == nil || update.Message.Chat.ID == chatId) {
			continue
		}

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, startStr)
			switch update.Message.Command() {				
			case "start":
				msg.ReplyMarkup = mainMenu1
			case "help":
				msg.Text = helpStr
			default:
				msg.Text = startStr
			}
			bot.Send(msg)
		}

		if (update.Message.Text == mainMenu1.Keyboard[0][0].Text || update.Message.Text == mainMenu1.Keyboard[0][1].Text) {

			probsummarySingMap[update.Message.From.ID] = new(probsummarySing)
			probsummarySingMap[update.Message.From.ID].State = 0
			
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			if update.Message.Text == mainMenu1.Keyboard[0][0].Text {
				probsummarySingMap[update.Message.From.ID].Category = "#сдам"
				msg.Text = "Введите текст вашего объявления. Не забудьте указать адрес, станцию метро и стоимость."
			}
			if update.Message.Text == mainMenu1.Keyboard[0][1].Text {
				probsummarySingMap[update.Message.From.ID].Category = "#сниму"
				msg.Text = "Введите текст вашего объявления. Не забудьте рассказать о себе. При необходимости укажите станцию метро, количество комнат."
			}
			msg.ReplyMarkup = mainMenuBack;
			bot.Send(msg)
			
		} else {
			cs, ok := probsummarySingMap[update.Message.From.ID]
			if (ok && (update.Message.Text != mainMenuBack.Keyboard[0][0].Text)) {
				if cs.State == 0 {
					cs.Text = update.Message.Text
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					if cs.Category == "#сдам"{
						msg.Text = "Добавьте фото квартиры."
					}
					if cs.Category == "#сниму"{
						msg.Text = "Добавьте ваше фото. После добавления нажмите \"Готово\""
					}
					msg.ReplyMarkup = mainMenuBackAndSkip
					bot.Send(msg)
					cs.State = 1
				} else if cs.State >= 1 {
					if update.Message.Photo!=nil {
						photo :=*update.Message.Photo
						cs.Photo = append(cs.Photo, photo[1].FileID)

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Изображение загружено. Нажмите \"Готово\" для завершения.")
						msg.ReplyMarkup = mainMenuComplete;
						bot.Send(msg)

						cs.State = cs.State + 1
					} else if (update.Message.Text == mainMenuComplete.Keyboard[0][0].Text || update.Message.Text == mainMenuBackAndSkip.Keyboard[0][1].Text){
						bot.Send(tgbotapi.NewMessage(chatId, cs.Category+"\n"+cs.Text+"\n\nДобавлено пользователем @"+update.Message.From.UserName+" через @"+bot.Self.UserName))
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
						msg.Text = "Ваше объявление будет опубликовано после проверки модератором на @channel_name\n\nТекст вашего объявления:\n\n"+cs.Text
						msg.ReplyMarkup = mainMenu1		
						bot.Send(msg)				
						for i := 0; i < len(cs.Photo); i++ {
							bot.Send(tgbotapi.NewPhotoShare(chatId, cs.Photo[i]))
							bot.Send(tgbotapi.NewPhotoShare(update.Message.Chat.ID, cs.Photo[i]))
						}	
						
						delete(probsummarySingMap, update.Message.From.ID)
					}
				}
			}
		}

		if update.Message.Text == mainMenuBack.Keyboard[0][0].Text {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, startStr)
			msg.ReplyMarkup = mainMenu1;
			bot.Send(msg)
			delete(probsummarySingMap, update.Message.From.ID)
		}

	}
}