package main

import (
	"log"

	telegrambot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

func main() {
	bot, err := telegramBot()
	if err != nil {
		log.Fatal(err)
	}

	db, err := dbConn()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	u := telegrambot.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		handleUpdate(bot, db, update)
	}
}
