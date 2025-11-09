package notifier

import (
	"context"
	"testing"
)

func TestLogSender(t *testing.T) {
	if err := (LogSender{}).Send(context.Background(), 42, "test message"); err != nil {
		t.Fatal(err)
	}
}
