package main

import (
	"log"
	"telebot-invent/api"
	"telebot-invent/conn"

	telegrambot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

func main() {
	bot, err := conn.TelegramBot()
	if err != nil {
		log.Fatal(err)
	}

	db, err := conn.DbConn()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	u := telegrambot.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		api.HandleUpdate(bot, db, update)
	}
}
