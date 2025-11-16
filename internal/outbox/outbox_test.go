package outbox

import (
	"errors"
	"testing"

	"github.com/arenkhachaturian/market-watchdog/internal/core"
)

func TestOutboxRetries(t *testing.T) {
	o := New(3)

	alert := core.AlertRule{ID: 1, Coin: "bitcoin"}
	o.Push(alert)

	sendFunc := func(alert core.AlertRule) error {
		// fail twice (attempts goes 1 -> 2), then succeed
		if o.attempts[alert.ID] < 2 {
			return errors.New("fail")
		}
		return nil
	}

	o.Process(sendFunc) // fail #1
	if o.attempts[alert.ID] != 1 { t.Fatalf("want attempts=1, got %d", o.attempts[alert.ID]) }

	o.Process(sendFunc) // fail #2
	if o.attempts[alert.ID] != 2 { t.Fatalf("want attempts=2, got %d", o.attempts[alert.ID]) }

	o.Process(sendFunc) // success; queue empty; attempts key removed
	if _, ok := o.attempts[alert.ID]; ok { t.Fatal("attempts key should be gone after success") }
	if len(o.PopAll()) != 0 { t.Fatal("queue should be empty after success") }

}
