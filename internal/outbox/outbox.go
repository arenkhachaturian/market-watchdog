package outbox

import (
	"sync"

	"github.com/arenkhachaturian/market-watchdog/internal/core"
)

// Outbox stores triggered alerts and handles retries
type Outbox struct {
	mu 			sync.Mutex
	queue		[]core.AlertRule
	attempts  	map[int]int
	maxRetry    int
}

func New(maxRetry int) *Outbox {
	return &Outbox{
		queue: []core.AlertRule{},
		attempts: make(map[int]int),
		maxRetry: maxRetry,
	}
}

// Push adds a new alert to the outbox
func (o *Outbox) Push(alert core.AlertRule) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.queue = append(o.queue, alert)
	o.attempts[alert.ID] = 0
}

// PopAll returns and clears all queued alerts
func (o *Outbox) PopAll() []core.AlertRule {
	o.mu.Lock()
	defer o.mu.Unlock()
	alerts := o.queue
	o.queue = []core.AlertRule{}
	return alerts
}

func (o *Outbox) Process(sendFunc func(core.AlertRule) error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	remaining := make([]core.AlertRule, 0, len(o.queue))
	for _, alert := range o.queue {
		if err := sendFunc(alert); err != nil {
			o.attempts[alert.ID]++               // increment on error only
			if o.attempts[alert.ID] <= o.maxRetry {
				remaining = append(remaining, alert)
			}
		} else {
			delete(o.attempts, alert.ID)         // clean up on success
		}
	}
	o.queue = remaining
}