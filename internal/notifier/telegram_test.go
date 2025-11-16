//go:build telegram

package notifier

import (
	"context"
	"os"
	"testing"
)

func TestTelegramSender_Send(t *testing.T) {
	token := os.Getenv("TELEGRAM_TOKEN")
	chatIDEnv := os.Getenv("TELEGRAM_CHAT_ID")
	if token == "" || chatIDEnv == "" {
		t.Fatalf("set TELEGRAM_TOKEN and TELEGRAM_CHAT_ID to run")
	}
	chatID, err := strconv.Atoi(chatIDEnv)
	if err != nil {
		t.Fatalf("bad TELEGRAM_CHAT_ID: %v", err)
	}

	sender, err := NewTelegramSender(token)
	if err != nil {
		t.Fatal(err)
	}
	if err := sender.Send(context.Background(), chatID, "market-watchdog test"); err != nil {
		t.Fatal(err)
	}
}
