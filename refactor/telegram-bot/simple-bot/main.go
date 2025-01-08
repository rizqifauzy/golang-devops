package main

import (
	"log"
	"os"

	telegrambot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	tokenbot := os.Getenv("TELEGRAM_TOKEN_BOT")

	if tokenbot == "" {
		log.Fatal("Please run export TELEGRAM_TOKEN_BOT=token_bot")
	}
	// Initialize the TelegramBot
	bot, err := telegrambot.NewBotAPI(tokenbot)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := telegrambot.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			//msg := telegrambot.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg := telegrambot.NewMessage(update.Message.Chat.ID, "You're connected")
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}
