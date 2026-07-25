package handlers

import("log"
"github.com/go-telegram/bot"
"fmt"
"mybot/timeapi"
"context"
)
func sendTime(
	ctx context.Context,
	b *bot.Bot,
	chatID int64,
	timeData timeapi.TimeResponse,
) {

	msg := fmt.Sprintf(
		"🕒 Current Time\n\n"+
			"📅 Date: %s\n"+
			"📆 Day: %s\n"+
			"⏰ Time: %s:%s:%s\n"+
			"🌍 Timezone: %s",
		timeData.Date,
		timeData.Day,
		timeData.Hour,
		timeData.Minute,
		timeData.Second,
		timeData.Timezone,
	)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msg,
	})

	if err != nil {
		log.Println(err)
	}
}