package main

import (
	"context"
	"fmt"
	"log"

	"bhatsapp/config"
	"bhatsapp/whatsapp"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

func main() {
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbContainer, err := sqlstore.New(context.Background(), "sqlite3", "file:whatsapp.db?_foreign_keys=on", nil)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	waClient, err := whatsapp.NewWhatsAppClient(dbContainer)

	if err != nil {
		fmt.Println("error getting client:", err)
		return
	}

	ctx := context.Background()
	if err := waClient.Login(ctx); err != nil {
		log.Fatalf("Failed to login: %v", err)
	}

	waClient.GetJoinedGroups()

	//Test id wont work
	if err := waClient.SendMessage(ctx, cfg.GroupChatID, "YOOOOOHOOOOO!"); err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	log.Println("Message sent successfully!", cfg.GroupChatID)
}
