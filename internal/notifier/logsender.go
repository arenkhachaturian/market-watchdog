package notifier

import (
	"context"
	"log"
)

type LogSender struct{}

func (LogSender) Send(_ context.Context, userID int, text string) error {
	log.Printf("[user:%d] %s", userID, text)
	return nil
}