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
		payload := update.Message.CommandArguments()
		if payload == "" {
			bot.Send(telegrambot.NewMessage(update.Message.Chat.ID, "Please input a server name or IP address.\nExample: /info serverapp1"))
			return
		}

		info, err := getServerInfo(db, payload)
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
	case "start":
		bot.Send(telegrambot.NewMessage(update.Message.Chat.ID, "Usage: /info <server_name_or_IPOAM>"))
	default:
		text := "Type /start to continue"
		msg := telegrambot.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyMarkup = CmdKeyboard()
		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}

}
