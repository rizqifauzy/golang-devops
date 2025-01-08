package main

import (
	"fmt"
	"log"

	telegrambot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

func telegramBot() (*telegrambot.BotAPI, error) {
	token := config("TELEGRAM_APITOKEN")
	bot, err := telegrambot.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("telegram API error: %s", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot, nil
}

func CmdKeyboard() telegrambot.ReplyKeyboardMarkup {
	var cmdKeyboard = telegrambot.NewReplyKeyboard(
		telegrambot.NewKeyboardButtonRow(
			telegrambot.NewKeyboardButton("/start"),
		),
		telegrambot.NewKeyboardButtonRow(
			telegrambot.NewKeyboardButton("/info"),
		),
	)
	return cmdKeyboard
}
