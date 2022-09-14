package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

const (
	helpMessage = `Вы ввели ключевое слово "Помощь"
Для упрощения коммуникаций в чат добавлены ключевые слова
для вызова ответа на часто задаваемые вопросы:
- Помощь
- Контакты
`
	contactsMessage = `Вы ввели ключевое слово "Контакты"
Контакты администратора:
- mail: yokozuna@yandex.ru
- tg:	@kkhitsko
`
)

const (
	replyInPrivateChannel = true
)

func getEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load("key.env")

	if err != nil {
		log.Fatalf("Error loading key.env file")
	}

	return os.Getenv(key)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(getEnvVariable("TG_API_KEY"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // Есть новое сообщение
			text := update.Message.Text      // Текст сообщения
			chatID := update.Message.Chat.ID //  ID чата
			userID := update.Message.From.ID // ID пользователя
			var replyMsg string

			log.Printf("[%s](%d) %s", update.Message.From.UserName, userID, text)

			// Анализируем текс сообщения
			switch {
			case text == "Помощь":
				replyMsg = helpMessage
				break
			case text == "Контакты":
				replyMsg = contactsMessage
				break
			}

			// В случае если сформирован автооответ - отправляем его
			if len(replyMsg) > 0 {
				if replyInPrivateChannel {
					chatID = userID
					msg := tgbotapi.NewMessage(chatID, replyMsg) // Создаем новое сообщение

					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(chatID, replyMsg)    // Создаем новое сообщение
					msg.ReplyToMessageID = update.Message.MessageID // Указываем сообщение, на которое нужно ответить

					bot.Send(msg)
				}
			}
		}
	}
}
