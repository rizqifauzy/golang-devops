package services

import (
	"fmt"
	"telegram-inventory/repositories"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	message := "Halo! Selamat datang di Bot Inventory Server. Gunakan perintah berikut untuk memulai:\n\n"
	message += "/info <nama_server/ipoam> - Untuk mendapatkan informasi server."
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	bot.Send(msg)
}
func GetServerInfo(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	keyword := update.Message.CommandArguments()
	if keyword == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Format salah! Gunakan /info <nama_server/ipoam>")
		bot.Send(msg)
		return
	}

	server, err := repositories.GetServerInfo(keyword)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Server tidak ditemukan.")
		bot.Send(msg)
		return
	}

	output := fmt.Sprintf(`Informasi Server:
VM: %s
IPOAM: %s
IPService: %s
Powerstate: %s
Datacenter: %s
OS according to the configuration file: %s
Apps Name: %s
Apps Priority: %s
Apps Custody (Email Address): %s
Prod/Non Prod: %s
Environment: %s
Site: %s
Managed By: %s
Support Level: %s
Notes: %s`,
		server.VMName,
		server.IPOAM,
		server.IPService,
		server.Powerstate,
		server.Datacenter,
		server.OSConfiguration,
		server.AppsName,
		server.AppsPriority,
		server.AppsCustody,
		server.ProdNonProd,
		server.Environment,
		server.Site,
		server.ManagedBy,
		server.SupportLevel,
		server.Notes,
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, output)
	bot.Send(msg)
}
