package handlers

import (
	"context"

	"log"
	"mybot/weather"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func MessageHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	// =========================
	// Handle Callback Queries
	// =========================

	if update.CallbackQuery != nil {

		data := update.CallbackQuery.Data

		switch data {

		case "weather":

			_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
			})
			if err != nil {
				log.Println(err)
			}

			/*if err != nil {
				_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.CallbackQuery.Message.Message.Chat.ID,
					Text:   "❌ Failed to fetch weather.",
				})
				return
			}*/

			chatID := update.CallbackQuery.Message.Message.Chat.ID

			waitingForCity[chatID] = true

			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "🌍 Please enter a city name:",
			})

			if err != nil {
				log.Println(err)
			}

		case "help":

			_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
			})
			if err != nil {
				log.Println(err)
			}

			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.Message.Message.Chat.ID,
				Text: `Available Commands:

/start
/help
/ping`,
			})

			if err != nil {
				log.Println(err)
			}
		}

		return
	}

	// Ignore updates without a message
	if update.Message == nil {
		return
	}

	text := update.Message.Text

	log.Println("Received:", text)

	if waitingForCity[update.Message.Chat.ID] {

		delete(waitingForCity, update.Message.Chat.ID)

		city := strings.TrimSpace(text)

		weatherData, err := weather.GetCurrent(city)
		if err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Failed to fetch weather.",
			})
			return
		}

		if len(weatherData.Weather) == 0 {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Weather information unavailable.",
			})
			return
		}

		sendWeather(ctx, b, update.Message.Chat.ID, weatherData)
		return
	}
	// =========================
	// /weather (no city)
	// =========================
	if text == "/weather" {

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Please provide a city.\n\nExample:\n/weather Delhi",
		})

		if err != nil {
			log.Println(err)
		}

		return
	}
	// =========================
	// /weather <city>
	// =========================
	if strings.HasPrefix(text, "/weather ") {

		city := strings.TrimSpace(strings.TrimPrefix(text, "/weather "))

		weatherData, err := weather.GetCurrent(city)
		if err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Failed to fetch weather.",
			})
			return
		}

		if len(weatherData.Weather) == 0 {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Weather information unavailable ❌",
			})
			return
		}
		sendWeather(ctx, b, update.Message.Chat.ID, weatherData)
		return
	}

	// =========================
	// /start
	// =========================
	if text == "/start" {

		keyboard := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text:         "🌤 Weather",
						CallbackData: "weather",
					},
				},
				{
					{
						Text:         "❓ Help",
						CallbackData: "help",
					},
				},
			},
		}

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "🤖 Welcome to my Go Telegram Bot!",
			ReplyMarkup: keyboard,
		})

		if err != nil {
			log.Println(err)
		}

		return
	}

	// =========================
	// /cancel
	// =========================
	if text == "/cancel" {

		if waitingForCity[update.Message.Chat.ID] {

			delete(waitingForCity, update.Message.Chat.ID)

			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Weather request cancelled.",
			})

			if err != nil {
				log.Println(err)
			}

			return
		}
	}

	// =========================
	// /help
	// =========================
	if text == "/help" {

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: `Available Commands:

/start
/help
/ping`,
		})

		if err != nil {
			log.Println(err)
		}

		return
	}

	// =========================
	// /ping
	// =========================
	if text == "/ping" {

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "🏓 Pong!",
		})

		if err != nil {
			log.Println(err)
		}

		return
	}

	// =========================
	// Echo
	// =========================
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})

	if err != nil {
		log.Println(err)
	}
}
