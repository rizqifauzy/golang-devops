package main

import (
	"database/sql"
	"log"

	telegrambot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

func handleUpdate(bot *telegrambot.BotAPI, db *sql.DB, update telegrambot.Update) {
	if update.Message == nil {
		return
	}

	if !isUserAuthorized(update.Message.From.ID) {
		bot.Send(telegrambot.NewMessage(update.Message.Chat.ID, "You are not authorized to use this bot."))
		return
	}

	if update.Message.IsCommand() {
		handleCommand(bot, db, update)
	}
}

func handleCommand(bot *telegrambot.BotAPI, db *sql.DB, update telegrambot.Update) {
	switch update.Message.Command() {
	case "info":
		query := update.Message.CommandArguments()
		if query == "" {
			bot.Send(telegrambot.NewMessage(update.Message.Chat.ID, "Usage: /info <server_name_or_IPOAM>"))
			return
		}

		info, err := getServerInfo(db, query)
		if err != nil {
			bot.Send(telegrambot.NewMessage(update.Message.Chat.ID, "Error retrieving server info."))
			log.Println("Error querying database:", err)
			return
		}

		if info == nil {
			bot.Send(telegrambot.NewMessage(update.Message.Chat.ID, "Server not found."))
			return
		}

		response := formatServerInfo(info)
		bot.Send(telegrambot.NewMessage(update.Message.Chat.ID, response))
	default:
		bot.Send(telegrambot.NewMessage(update.Message.Chat.ID, "Unknown command."))
	}
}
