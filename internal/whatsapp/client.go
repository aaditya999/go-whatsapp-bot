package whatsapp

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aaditya999/go-whatsapp-bot/internal/config"
	"go.mau.fi/whatsmeow"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

var processedCount int

type WhatsAppClient struct {
	client         *whatsmeow.Client
	eventHandlerID uint32
	config         *config.Config
}

func getConfig() *config.Config {
	cfg, err := config.LoadConfig("internal/config/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return cfg
}

func NewWhatsAppClient(dbContainer *sqlstore.Container) (*WhatsAppClient, error) {
	deviceStore, err := dbContainer.GetFirstDevice(context.Background())
	if err != nil {
		return nil, err
	}
	client := whatsmeow.NewClient(deviceStore, nil)
	cfg := getConfig()
	return &WhatsAppClient{client: client, config: cfg}, nil
}

func (wac *WhatsAppClient) Login(ctx context.Context) error {
	if wac.client.Store.ID == nil {
		qrChan, _ := wac.client.GetQRChannel(ctx)
		if err := wac.client.Connect(); err != nil {
			return err
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Printf("Scan this QR code with your WhatsApp app: %s\n", evt.Code)
			} else if evt.Event == "success" || evt.Event == "timeout" || evt.Event == "error" {
				break
			}
		}
	} else {
		return wac.client.Connect()
	}
	return nil
}

func (wac *WhatsAppClient) Logout(ctx context.Context) {
	wac.client.Disconnect()
}

func (wac *WhatsAppClient) GetJoinedGroups() {
	groups, err := wac.client.GetJoinedGroups()
	if err != nil {
		fmt.Println("Error fetching groups:", err)
		return
	}
	for _, group := range groups {
		fmt.Println("Group ID:", group.JID, "Group Name:", group.Name)
	}
}

func (wac *WhatsAppClient) SendMessage(ctx context.Context, groupID string, message string) error {
	jid := types.NewJID(groupID, "g.us")
	msg := &waE2E.Message{Conversation: proto.String(message)}
	resp, err := wac.client.SendMessage(context.Background(), jid, msg)
	if err != nil {
		fmt.Errorf("Error sending message: %v", err)
	} else {
		fmt.Printf("Message sent (server timestamp: %s)", resp.Timestamp)
	}
	return nil
}

func (wac *WhatsAppClient) Register() {
	wac.eventHandlerID = wac.client.AddEventHandler(wac.myEventHandler)
}

func (wac *WhatsAppClient) myEventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		processedCount++
		fmt.Printf("Processed %d messages. Latest timestamp: %v (%s)\n", processedCount, v.Info.Timestamp, v.Info.Timestamp.Format(time.RFC3339))
		// Only respond to messages in the target group
		chatJID := v.Info.MessageSource.Chat.String()
		if chatJID == wac.config.GroupChatID {
			conv := v.Message.GetConversation()
			fmt.Println("Received a message!", conv)
			if conv == "/speak" {
				// Respond with a message
				groupId := strings.TrimSuffix(chatJID, "@g.us")
				go wac.SendMessage(context.Background(), groupId, "Bhow bhow")
			}
		}
	}
}
