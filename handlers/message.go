package handlers

import (
	"context"
	"fmt"
	"log"
	forecast "mybot/Forecast"
	"mybot/database"
	"mybot/geocoding"
	"mybot/timeapi"
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

		city, exists := database.GetUserCity(update.Message.Chat.ID)

		if !exists {

			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ No default city set.\nUse:\n/setcity Delhi",
			})

			if err != nil {
				log.Println(err)
			}

			return
		}

		weatherData, err := weather.GetCurrent(city)
		if err != nil {

			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Failed to fetch weather.",
			})

			return
		}

		sendWeather(ctx, b, update.Message.Chat.ID, weatherData)
		return
	}
	// =========================
	// /time (no city)
	// =========================
	if text == "/time" {

		city, exists := database.GetUserCity(update.Message.Chat.ID)

		if !exists {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ No default city set.\nUse:\n/setcity Delhi",
			})
			return
		}

		cityInfo, err := geocoding.GetCity(city)
		if err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ City not found.",
			})
			return
		}

		timeData, err := timeapi.GetCurrentTime(cityInfo.Timezone)
		if err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Failed to fetch time.",
			})
			return
		}

		sendTime(ctx, b, update.Message.Chat.ID, timeData)
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
	// /forecast (no city)
	// =========================
	if text == "/forecast" {

		city, exists := database.GetUserCity(update.Message.Chat.ID)

		if !exists {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ No default city set.\nUse:\n/setcity Delhi",
			})
			return
		}

		forecastData, err := forecast.GetForecast(city)
		if err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Failed to fetch forecast.",
			})
			return
		}

		msg := fmt.Sprintf("📍 %s, %s\n\n🌤 Forecast\n\n",
			forecastData.City.Name,
			forecastData.City.Country,
		)

		for i, item := range forecastData.List {

			if i == 5 {
				break
			}

			msg += fmt.Sprintf(
				"🕒 %s\n🌡 %.1f°C\n☁ %s\n\n",
				item.DateTime,
				item.Main.Temp,
				item.Weather[0].Description,
			)
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   msg,
		})

		if err != nil {
			log.Println(err)
		}

		return
	}

	// =========================
	// /forecast <city>
	// =========================
	if strings.HasPrefix(text, "/forecast ") {

		city := strings.TrimSpace(strings.TrimPrefix(text, "/forecast "))

		forecastData, err := forecast.GetForecast(city)
		if err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Failed to fetch forecast.",
			})
			return
		}

		msg := fmt.Sprintf("📍 %s, %s\n\n🌤 Forecast\n\n",
			forecastData.City.Name,
			forecastData.City.Country,
		)

		for i, item := range forecastData.List {

			if i == 5 {
				break
			}

			msg += fmt.Sprintf(
				"🕒 %s\n🌡 %.1f°C\n☁ %s\n\n",
				item.DateTime,
				item.Main.Temp,
				item.Weather[0].Description,
			)
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   msg,
		})

		if err != nil {
			log.Println(err)
		}

		return
	}

	// =========================
	// /setcity
	// =========================
	if strings.HasPrefix(text, "/setcity ") {

		city := strings.TrimSpace(strings.TrimPrefix(text, "/setcity "))

		if city == "" {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Please provide a city.",
			})
			return
		}

		database.SaveUserCity(update.Message.Chat.ID, city)

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("✅ Default city set to %s.", city),
		})

		if err != nil {
			log.Println(err)
		}

		return
	}

	// =========================
	// /time <city>
	// =========================
	if strings.HasPrefix(text, "/time ") {

		city := strings.TrimSpace(strings.TrimPrefix(text, "/time "))

		cityInfo, err := geocoding.GetCity(city)
		if err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ City not found.",
			})
			return
		}

		timeData, err := timeapi.GetCurrentTime(cityInfo.Timezone)
		if err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Failed to fetch time.",
			})
			return
		}

		msg := fmt.Sprintf(
			"🕒 %s, %s\n\n🌍 Timezone: %s\n\n⏰ %s",
			cityInfo.Name,
			cityInfo.Country,
			cityInfo.Timezone,
			timeData.Datetime,
		)

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   msg,
		})

		if err != nil {
			log.Println(err)
		}

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
