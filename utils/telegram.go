package utils

import "github.com/go-telegram/bot"

func MakeMessageParams(t string, mi int) *bot.EditMessageTextParams {
	return &bot.EditMessageTextParams{
		ChatID:    "@pactus_status",
		Text:      t,
		MessageID: mi,
	}
}
