package main

import (
	"context"
	"log"
	"os"

	"mybot/database"
	"mybot/handlers"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables (for local development)
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using Railway environment variables")
	}
	// Read bot token
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is empty")
	}

	database.InitDatabase()

	// Create bot
	b, err := bot.New(token, bot.WithDefaultHandler(handlers.MessageHandler))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot is running...")

	// Start listening for updates
	b.Start(context.Background())
}
