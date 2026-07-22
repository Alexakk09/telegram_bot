package handlers

import (
	"context"
	"fmt"
	"log"
	"mybot/weather"

	"github.com/go-telegram/bot"
)

func sendWeather(ctx context.Context, b *bot.Bot, chatID int64, weatherData weather.WeatherResponse) {

	msg := fmt.Sprintf(
		"🌍 City: %s\n\n🌡 Temperature: %.1f°C\n💧 Humidity: %d%%\n☁ Weather: %s",
		weatherData.Name,
		weatherData.Main.Temp,
		weatherData.Main.Humidity,
		weatherData.Weather[0].Description,
	)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msg,
	})

	if err != nil {
		log.Println(err)
	}
}