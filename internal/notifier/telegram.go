package notifier

import (
	"context"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramSender struct {
	api *tgbot.BotAPI
}

func NewTelegramSender(botToken string) (*TelegramSender, error) {
	api, err := tgbot.NewBotAPI(botToken)
	if err  != nil {
		return nil, err
	}
	return &TelegramSender{api: api}, nil
}

func (telegram *TelegramSender) Send(_ context.Context, userID int, text string) error {
	msg := tgbot.NewMessage(int64(userID), text)
	_, err := telegram.api.Send(msg)
	return err
}