package main

import (
	"encoding/json"
	"log"
	"os"

	r "github.com/mimin0/go-budget-bot/spreadSheet"

	"gopkg.in/telegram-bot-api.v4"
)

type Config struct {
	TelegramBotToken string
	SpreadID         string
}

func main() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}
	// fmt.Println(configuration.TelegramBotToken)

	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}
	//new messages will be loaded into the chanel 'updates'
	for update := range updates {

		command := update.Message.Command()
		if command == "" {
			// condition for plane taxt messages. NonCommand
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		} else {
			// condition for command messages
			switch command {
			case "help":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				msg.Text = "/help - show help info \n" +
					"/showall - show all items in a budget spreadsheet \n" +
					"/add - adding new item. for instance: \n [/add] [short_comment] [summ] [type]  '/add beer 50 fun'\n" +
					"/showtypes - show all types"
				bot.Send(msg)
			case "showall":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, r.Reader(configuration.SpreadID))
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			case "add":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, r.Writer(configuration.SpreadID, update.Message.Date, update.Message.Text))
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			case "showtypes":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				msg.Text = "Groceries, Gifts, Appartment, Transportation, Departmental, Travel, Loan, OneTime, Car, Fun"
				bot.Send(msg)
			}
		}

	}

}
