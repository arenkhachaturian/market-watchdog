package notifier

import (
	"context"
)

// Sender abstracts how we deliver messages.
type Sender interface {
	Send(ctx context.Context, userID int, text string) error
}
