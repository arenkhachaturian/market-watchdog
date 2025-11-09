package core

import (
	"time"

	"github.com/arenkhachaturian/market-watchdog/internal/fetcher"
)

type Evaluator struct {
	Fetch fetcher.PriceFetcher
}

func NewEvaluator(f fetcher.PriceFetcher) *Evaluator {
	return &Evaluator{Fetch: f}
}

// EvaluateRules checks prices against rules and returns triggered alerts
func (e *Evaluator) EvaluateRules(rules []AlertRule, now time.Time) ([]AlertRule, error) {
	triggered := []AlertRule{}

	for i := range rules {
		rule := &rules[i]

		if !rule.CanTrigger(now) {
			continue
		}

		price, err := e.Fetch.GetPrice(rule.Coin)
		if err != nil {
			return nil, err
		}

		switch rule.Comparator {
		case Greater:
			if price > rule.Threshold {
				triggered = append(triggered, *rule)
				rule.LastTriggered = &now
			}
		case Less:
			if price < rule.Threshold {
				triggered = append(triggered, *rule)
				rule.LastTriggered = &now
			}
		}
	}
	return triggered, nil
}
