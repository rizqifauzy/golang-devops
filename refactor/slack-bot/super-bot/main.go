package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type ServerInfo struct {
	VM, IPOAM, IPService, Powerstate, Datacenter, OS, AppsName, AppsPriority, AppsCustody,
	ProdNonProd, Environment, Site, ManagedBy, SupportLevel, Notes string
}

func main() {
	loadEnv()

	client, socketClient := setupSlackClient()
	db := setupDatabase()
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handleSocketEvents(ctx, client, socketClient, db)
	socketClient.Run()
}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func setupSlackClient() (*slack.Client, *socketmode.Client) {
	client := slack.New(
		os.Getenv("SLACK_AUTH_TOKEN"),
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(os.Getenv("SLACK_APP_TOKEN")),
	)
	return client, socketmode.New(client)
}

func setupDatabase() *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("PG_HOST"), os.Getenv("PG_USER"), os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DB_NAME"), os.Getenv("PG_PORT"), os.Getenv("PG_SSL_MODE"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

func handleSocketEvents(ctx context.Context, client *slack.Client, socketClient *socketmode.Client, db *sql.DB) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down socketmode listener")
			return
		case event := <-socketClient.Events:
			handleEvent(event, client, socketClient, db)
		}
	}
}

func handleEvent(event socketmode.Event, client *slack.Client, socketClient *socketmode.Client, db *sql.DB) {
	switch event.Type {
	case socketmode.EventTypeInteractive:
		handleInteractiveEvent(event, socketClient)
	case socketmode.EventTypeSlashCommand:
		handleSlashCommand(event, client, socketClient, db)
	case socketmode.EventTypeEventsAPI:
		handleEventsAPI(event, client, socketClient)
	}
}

func handleInteractiveEvent(event socketmode.Event, socketClient *socketmode.Client) {
	callback, ok := event.Data.(slack.InteractionCallback)
	if !ok || callback.Type != slack.InteractionTypeBlockActions {
		log.Println("Ignored unsupported event")
		return
	}
	socketClient.Ack(*event.Request)
}

func handleSlashCommand(event socketmode.Event, client *slack.Client, socketClient *socketmode.Client, db *sql.DB) {
	command, ok := event.Data.(slack.SlashCommand)
	if !ok {
		log.Println("Ignored unsupported slash command")
		return
	}
	socketClient.Ack(*event.Request)
	processSlashCommand(command, client, db)
}

func processSlashCommand(command slack.SlashCommand, client *slack.Client, db *sql.DB) {
	user, _ := client.GetUserInfo(command.UserID)
	switch command.Command {
	case "/info":
		query := strings.TrimSpace(command.Text)
		if query == "" {
			postSlackMessage(client, command.ChannelID, "Usage: /info <server_name> or <IP_OAM>\n", "#f7b20f")
			return
		}
		info, err := getServerInfo(db, query)
		if err != nil {
			postSlackMessage(client, command.ChannelID, "Error retrieving server info.", "#f00202")
			log.Println("Database error:", err)
			return
		}
		if info == nil {
			postSlackMessage(client, command.ChannelID, fmt.Sprintf("Server %s not found.", query), "#0e7fe8")
		} else {
			response := formatServerInfo(info)
			postSlackMessage(client, command.ChannelID, response, "#4af030")
		}
	default:
		postSlackMessage(client, command.ChannelID, fmt.Sprintf("Hello %s! Use /info for server information.", user.Name), "#4af030")
	}
}

func postSlackMessage(client *slack.Client, channelID, text, color string) {
	attachment := slack.Attachment{
		Text:  text,
		Color: color,
		Fields: []slack.AttachmentField{
			{
				Title: "Date",
				Value: time.Now().String(),
			},
		},
	}
	client.PostMessage(channelID, slack.MsgOptionAttachments(attachment))
}

func handleEventsAPI(event socketmode.Event, client *slack.Client, socketClient *socketmode.Client) {
	eventsAPI, ok := event.Data.(slackevents.EventsAPIEvent)
	if !ok {
		log.Printf("Unsupported EventsAPIEvent: %v", event)
		return
	}
	socketClient.Ack(*event.Request)
	if err := handleAppMention(eventsAPI, client); err != nil {
		log.Println("Error handling event:", err)
	}
}

func handleAppMention(event slackevents.EventsAPIEvent, client *slack.Client) error {
	if innerEvent, ok := event.InnerEvent.Data.(*slackevents.AppMentionEvent); ok {
		return respondToAppMention(innerEvent, client)
	}
	return errors.New("unsupported event type")
}

func respondToAppMention(event *slackevents.AppMentionEvent, client *slack.Client) error {
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}
	text := strings.ToLower(event.Text)
	response := "I am good. How are you?"
	if strings.Contains(text, "hello") || strings.Contains(text, "hi") {
		response = fmt.Sprintf("Hello %s!", user.Name)
	} else if strings.Contains(text, "weather") {
		response = fmt.Sprintf("The weather is sunny today, %s!", user.Name)
	}
	postSlackMessage(client, event.Channel, response, "#4af030")
	return nil
}

func getServerInfo(db *sql.DB, query string) (*ServerInfo, error) {
	var info ServerInfo
	err := db.QueryRow(`SELECT * FROM servers WHERE vm_name = $1 OR ip_oam = $1`, query).Scan(
		&info.VM, &info.IPOAM, &info.IPService, &info.Powerstate, &info.Datacenter, &info.OS,
		&info.AppsName, &info.AppsPriority, &info.AppsCustody, &info.ProdNonProd, &info.Environment,
		&info.Site, &info.ManagedBy, &info.SupportLevel, &info.Notes,
	)
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
OS: %s
Apps Name: %s
Apps Priority: %s
Apps Custody: %s
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
