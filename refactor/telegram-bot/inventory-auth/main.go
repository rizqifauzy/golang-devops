package main

import (
	"log"

	telegrambot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize bot
	bot, err := telegramBot()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize database connection
	db, err := dbConn()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Start listening for updates
	u := telegrambot.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // Ignore non-Message updates
			continue
		}

		if !isUserAuthorized(update.Message.From.ID) {
			bot.Send(telegrambot.NewMessage(update.Message.Chat.ID, "You are not authorized to use this bot."))
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "info":
				query := update.Message.CommandArguments()
				if query == "" {
					bot.Send(telegrambot.NewMessage(update.Message.Chat.ID, "Usage: /info <server_name_or_IPOAM>"))
					continue
				}

				info, err := getServerInfo(db, query)
				if err != nil {
					bot.Send(telegrambot.NewMessage(update.Message.Chat.ID, "Error retrieving server info."))
					log.Println("Error querying database:", err)
					continue
				}

				if info == nil {
					bot.Send(telegrambot.NewMessage(update.Message.Chat.ID, "Server not found."))
				} else {
					response := formatServerInfo(info)
					bot.Send(telegrambot.NewMessage(update.Message.Chat.ID, response))
				}
			default:
				bot.Send(telegrambot.NewMessage(update.Message.Chat.ID, "Unknown command."))
			}
		}
	}
}
