package whatsapp

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

type WhatsAppClient struct {
	client *whatsmeow.Client
}

func NewWhatsAppClient(dbContainer *sqlstore.Container) (*WhatsAppClient, error) {
	deviceStore, err := dbContainer.GetFirstDevice(context.Background())
	if err != nil {
		return nil, err
	}
	client := whatsmeow.NewClient(deviceStore, nil)
	return &WhatsAppClient{client: client}, nil
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
		fmt.Print("Message sent (server timestamp: %s)", resp.Timestamp)
	}
	return nil
}
