package main

import (
	"context"
	"log"
	"os"

	"mybot/handlers"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
	"mybot/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
	log.Fatalf("godotenv.Load failed: %v", err)
}

	// Read bot token
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is empty")
	}

	database.InitDatabase()
	
	// Create bot
	b, err := bot.New(token, bot.WithDefaultHandler(handlers.MessageHandler),)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot is running...")

	// Start listening for updates
	b.Start(context.Background())
}

