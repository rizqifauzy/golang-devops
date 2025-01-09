package main

import (
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

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Slack client setup
	token := os.Getenv("SLACK_AUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")
	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))
	socketClient := socketmode.New(client)

	// Database connection setup
	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("PG_HOST"), os.Getenv("PG_USER"), os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DB_NAME"), os.Getenv("PG_PORT"), os.Getenv("PG_SSL_MODE"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Goroutine to handle incoming socket events
	go func() {
		for event := range socketClient.Events {
			switch event.Type {
			case socketmode.EventTypeInteractive:
				callback, ok := event.Data.(slack.InteractionCallback)
				if !ok {
					log.Println("Ignored unsupported event")
					continue
				}
				if callback.Type == slack.InteractionTypeBlockActions {
					socketClient.Ack(*event.Request)
					continue
				}

			case socketmode.EventTypeSlashCommand:

				command, ok := event.Data.(slack.SlashCommand)
				if !ok {
					log.Println("Ignored unsupported slash command")
					continue
				}
				socketClient.Ack(*event.Request)
				handleSlashCommand(command, client, db)
			case socketmode.EventTypeEventsAPI:
				eventsAPI, ok := event.Data.(slackevents.EventsAPIEvent)
				if !ok {
					log.Printf("Could not typecast the event to EventsAPIEvent: %v\n", event)
					continue
				}

				socketClient.Ack(*event.Request)

				err := HandleEventMessage(eventsAPI, client)
				if err != nil {
					log.Printf("Error handling event: %v", err)
				}
			}
		}
	}()

	// Start the socket client
	socketClient.Run()
}

func handleSlashCommand(command slack.SlashCommand, client *slack.Client, db *sql.DB) {
	user, err := client.GetUserInfo(command.UserID)
	if err != nil {
		client.PostMessage(command.ChannelID, slack.MsgOptionText("error get user info", false))
		return
	}
	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			{
				Title: "Date",
				Value: time.Now().String(),
			},
		},
	}

	switch command.Command {
	case "/info":
		query := command.Text

		if query == "" {
			newAttachment := slack.Attachment{
				Pretext: "Please input the name or IP of the server:",
				Text:    "Usage: /info <server_name> or <IP_OAM>\n",
				Color:   "#f7b20f",
				Fields:  attachment.Fields,
			}

			client.PostMessage(command.ChannelID, slack.MsgOptionAttachments(newAttachment))
			return
		}

		info, err := getServerInfo(db, query)
		if err != nil {
			newAttachment := slack.Attachment{
				Pretext: "Got Error :",
				Text:    "Error retrieving server info.\n",
				Color:   "##f00202",
				Fields:  attachment.Fields,
			}
			client.PostMessage(command.ChannelID, slack.MsgOptionAttachments(newAttachment))
			log.Println("Error querying database:", err)
			return
		}

		if info == nil {
			newAttachment := slack.Attachment{
				Pretext: "Not Found :",
				Text:    fmt.Sprintf("Server %s not found.\n", query),
				Color:   "##0e7fe8",
				Fields:  attachment.Fields,
			}
			client.PostMessage(command.ChannelID, slack.MsgOptionAttachments(newAttachment))
		} else {
			response := formatServerInfo(info)
			newAttachment := slack.Attachment{
				Pretext: "Here is the result :",
				Text:    response,
				Color:   "#4af030",
				Fields:  attachment.Fields,
			}

			client.PostMessage(command.ChannelID, slack.MsgOptionAttachments(newAttachment))

		}

	default:
		newAttachment := slack.Attachment{
			Pretext: fmt.Sprintf("Hello %s how are you today.", user.Name),
			Text:    "For get the server information, use comamand /info\n",
			Color:   "#4af030",
			Fields:  attachment.Fields,
		}
		client.PostMessage(command.ChannelID, slack.MsgOptionAttachments(newAttachment))
	}
}

func HandleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
	switch event.Type {
	case slackevents.CallbackEvent:
		innerEvent := event.InnerEvent
		switch evnt := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			err := HandleAppMentionEventToBot(evnt, client)
			if err != nil {
				return err
			}
		default:
			return errors.New("unsupported event type")
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}

func HandleAppMentionEventToBot(event *slackevents.AppMentionEvent, client *slack.Client) error {
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	text := strings.ToLower(event.Text)

	attachment := slack.Attachment{}

	if strings.Contains(text, "hello") || strings.Contains(text, "hi") {
		attachment.Text = fmt.Sprintf("Hello %s!", user.Name)
		attachment.Color = "#4af030"
	} else if strings.Contains(text, "weather") {
		attachment.Text = fmt.Sprintf("The weather is sunny today, %s!", user.Name)
		attachment.Color = "#4af030"
	} else {
		attachment.Text = fmt.Sprintf("I am good. How are you, %s?", user.Name)
		attachment.Color = "#4af030"
	}

	_, _, err = client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}

	return nil
}

func getServerInfo(db *sql.DB, query string) (*ServerInfo, error) {
	var info ServerInfo
	row := db.QueryRow(`SELECT * FROM servers WHERE vm_name = $1 OR ip_oam = $1`, query)

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
