package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type ServerInfo struct {
	VM           string
	IPOAM        string
	IPService    string
	Powerstate   string
	Datacenter   string
	OS           string
	AppsName     string
	AppsPriority string
	AppsCustody  string
	ProdNonProd  string
	Environment  string
	Site         string
	ManagedBy    string
	SupportLevel string
	Notes        string
}

func config(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}

func getServerInfo(db *sql.DB, query string) (*ServerInfo, error) {
	var info ServerInfo
	row := db.QueryRow(`SELECT vm_name, ip_oam, ip_service, powerstate, datacenter, os_configuration, apps_name, 
		apps_priority, apps_custody, prod_non_prod, environment, site, managed_by, 
		support_level, notes FROM servers WHERE vm_name = $1 OR ip_oam = $1`, query)

	err := row.Scan(&info.VM, &info.IPOAM, &info.IPService, &info.Powerstate, &info.Datacenter, &info.OS,
		&info.AppsName, &info.AppsPriority, &info.AppsCustody, &info.ProdNonProd, &info.Environment,
		&info.Site, &info.ManagedBy, &info.SupportLevel, &info.Notes)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &info, nil
}

func formatServerInfo(info *ServerInfo) string {
	return fmt.Sprintf(`Informasi Server:
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
		info.VM, info.IPOAM, info.IPService, info.Powerstate, info.Datacenter, info.OS,
		info.AppsName, info.AppsPriority, info.AppsCustody, info.ProdNonProd, info.Environment,
		info.Site, info.ManagedBy, info.SupportLevel, info.Notes)
}

func main() {
	// Initialize bot
	token := config("TELEGRAM_APITOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Initialize database connection
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config("PG_HOST"), config("PG_USER"), config("PG_PASSWORD"), config("PG_DB_NAME"), config("PG_PORT"), config("PG_SSL_MODE"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Start listening for updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // Ignore non-Message updates
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "info":
				query := update.Message.CommandArguments()
				if query == "" {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: /info <server_name_or_IPOAM>"))
					continue
				}

				info, err := getServerInfo(db, query)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error retrieving server info."))
					log.Println("Error querying database:", err)
					continue
				}

				if info == nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Server not found."))
				} else {
					response := formatServerInfo(info)
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, response))
				}
			default:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command."))
			}
		}
	}
}
