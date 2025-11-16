package core

import (
	"context"

	"github.com/arenkhachaturian/market-watchdog/internal/notifier"
)

type outboxPusher interface{ Push(AlertRule) }
type outboxProc interface {
	Process(func(AlertRule) error)
	PopAll() []AlertRule
}

// EnqueueMatches pushes evaluator matches into the outbox.
func EnqueueMatches(ob outboxPusher, matches []AlertRule) {
	for _, m := range matches {
		ob.Push(m)
	}
}

// DeliverOnce processes the outbox once using the provided Sender.
func DeliverOnce(ctx context.Context, ob outboxProc, sender notifier.Sender) {
	ob.Process(func(a AlertRule) error {
		return sender.Send(ctx, a.UserID, FormatAlert(a))
	})
}