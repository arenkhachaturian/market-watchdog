// before:
// package core
// after:
package core_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	corepkg "github.com/arenkhachaturian/market-watchdog/internal/core"
	"github.com/arenkhachaturian/market-watchdog/internal/fetcher"
	"github.com/arenkhachaturian/market-watchdog/internal/outbox"
)

type captureSender struct {
	mu   sync.Mutex
	sent []string
}

func (c *captureSender) Send(_ context.Context, _ int, text string) error {
	c.mu.Lock(); defer c.mu.Unlock()
	c.sent = append(c.sent, text)
	return nil
}

func TestEndToEndFlow(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]map[string]float64{"bitcoin": {"usd": 70000}})
	}))
	defer s.Close()

	f := fetcher.NewCoinGeckoFetcher()
	f.BaseURL = s.URL

	e := corepkg.NewEvaluator(f)
	ob := outbox.New(3)
	cs := &captureSender{}

	rules := []corepkg.AlertRule{
		{ID: 1, UserID: 99, Coin: "bitcoin", Comparator: corepkg.Greater, Threshold: 50000, CooldownMin: 10},
	}

	now := time.Now()
	matches, err := e.EvaluateRules(rules, now)
	if err != nil { t.Fatal(err) }
	if len(matches) != 1 { t.Fatalf("expected 1 match, got %d", len(matches)) }

	corepkg.EnqueueMatches(ob, matches)
	corepkg.DeliverOnce(context.Background(), ob, cs)

	if got := len(ob.PopAll()); got != 0 { t.Fatalf("outbox not empty, got %d", got) }
	if len(cs.sent) != 1 { t.Fatalf("expected 1 send, got %d", len(cs.sent)) }

	matches2, _ := e.EvaluateRules(rules, now)
	if len(matches2) != 0 { t.Fatalf("cooldown failed, got %d", len(matches2)) }
}
