package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aaditya999/go-whatsapp-bot/internal/whatsapp"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

func main() {

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

	waClient.Register()
	fmt.Println("WhatsApp client is running...")
	// if err := waClient.SendMessage(ctx, cfg.GroupChatID, "YOOOOOHOOOOO!"); err != nil {
	// 	log.Fatalf("Failed to send message: %v", err)
	// }

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("Shutting down WhatsApp client...")
	waClient.Logout(ctx)
}
