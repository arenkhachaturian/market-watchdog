package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/arenkhachaturian/market-watchdog/internal/core"
	"github.com/arenkhachaturian/market-watchdog/internal/store"
)

type TelegramBot struct {
	api 	  *tgbot.BotAPI
	alertRepo store.alertRepo
}

func NewTelegramBot(botToken string, alertRepo store.AlertRepo) (*TelegramBot, error) {
	api, err := tgbot.NewBotApi(botToken)
	if err != nil {
		return nil, err
	}
	return &TelegramBot{api: api, alertRepo: alertRepo}, nil
}

func (telegramBot *TelegramBot) Run(ctx context.Context) error {
	updateConfig := tgbot.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := telegramBot.api.GetUpdatesChan(updateConfig)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			telegramBot.handle(update.Message)
		}
	}
}

func (telegramBot *TelegramBot) handle(message *tgbot.Message) {
	text := strings.TrimSpace(message.Text)
	switch {
	case strings.HasPrefix(text, "/start"):
		telegramBot.reply(message.Chat.ID, "Commands: /add <coin> <gt|lt> <value>, /list, /rm <id>")
	case strings.HasPrefix(text, "/add")
		telegram.cmdAdd(message)
	case strigns.HasPrefix(text, "/list")
		telegram.cmdList(message)
	case strings.HasPrefix(text, "/rm"):
		telegramBot.cmdRemove(message)
	default:
		telegramBot.reply(message.Chat.ID, "Try: /add <coin> <gt|lt> <value>")
	}
}

func (telegramBot *TelegramBot) cmdAdd(message *tgbot.Message) {
	fields := strings.Fields(message.Text)
	if len(fields) != 4 {
		telegramBot.reply(message.Chat.ID, "Usage: /add <coin> <gt|lt> <value>")
		return
	}
	coin := strings.ToLower(fields[1])

	var comparator core.Comparator
	switch fields[2] {
	case "gt":
		comparator = core.Greater
	case "lt":
		comparator = core.Less
	default:
		telegramBot.reply(message.Chat.ID, "Operator must be gt or lt")
		return
	}

	threshold, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		telegramBot.reply(message.Chat.ID, "Value must be a number")
		return
	}

	newID, err := telegramBot.alertRepo.Create(message.Context(), core.AlertRule{
		UserID:      int(message.Chat.ID), // treat chat ID as user ID for now
		Coin:        coin,
		Comparator:  comparator,
		Threshold:   threshold,
		CooldownMin: 30,
	})
	if err != nil {
		log.Println("create alert:", err)
		telegramBot.reply(message.Chat.ID, "Failed to save")
		return
	}
	telegramBot.reply(message.Chat.ID, fmt.Sprintf("Added #%d: %s %s %.4f", newID, coin, fields[2], threshold))
}

func (telegramBot *TelegramBot) cmdList(message *tgbot.Message) {
	alerts, err := telegramBot.alertRepo.ListByUser(message.Context(), int(message.Chat.ID))
	if err != nil {
		telegramBot.reply(message.Chat.ID, "Failed to list")
		return
	}
	if len(alerts) == 0 {
		telegramBot.reply(message.Chat.ID, "No alerts")
		return
	}
	var builder strings.Builder
	for _, alert := range alerts {
		operator := map[core.Comparator]string{core.Greater: "gt", core.Less: "lt"}[alert.Comparator]
		fmt.Fprintf(&builder, "#%d %s %s %.4f cd=%d\n",
			alert.ID, alert.Coin, operator, alert.Threshold, alert.CooldownMin)
	}
	telegramBot.reply(message.Chat.ID, builder.String())
}

func (telegramBot *TelegramBot) cmdRemove(message *tgbot.Message) {
	fields := strings.Fields(message.Text)
	if len(fields) != 2 {
		telegramBot.reply(message.Chat.ID, "Usage: /rm <id>")
		return
	}
	alertID, err := strconv.Atoi(fields[1])
	if err != nil {
		telegramBot.reply(message.Chat.ID, "Id must be a number")
		return
	}
	if err := telegramBot.alertRepo.Delete(message.Context(), int(message.Chat.ID), alertID); err != nil {
		telegramBot.reply(message.Chat.ID, "Not found")
		return
	}
	telegramBot.reply(message.Chat.ID, "Removed")
}

func (telegramBot *TelegramBot) reply(chatID int64, text string) {
	_, _ = telegramBot.api.Send(tgbot.NewMessage(chatID, text))
}