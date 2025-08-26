package config

import (
	"os"
	"encoding/json"
)

type Config struct {
	WhatsAppNumber string `json:"whatsapp_number"`
	GroupChatID    string `json:"group_chat_id"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}