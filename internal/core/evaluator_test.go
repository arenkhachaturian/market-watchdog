package core

import (
	"testing"
	"time"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"github.com/arenkhachaturian/market-watchdog/internal/fetcher"
)

func TestEvaluator(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]map[string]float64{"bitcoin": {"usd": 70000}})
	}))
	defer mock.Close()
	
	f := fetcher.NewCoinGeckoFetcher()
	f.BaseURL = mock.URL 


	rules := []AlertRule{
		{ID: 1, UserID: 42, Coin: "bitcoin", Comparator: Greater, Threshold: 50000, CooldownMin: 10},
	}

	eval := NewEvaluator(f)

	// fake now
	now := time.Now()
	triggered, err := eval.EvaluateRules(rules, now)
	if err != nil {
		t.Fatal(err)
	}

	if len(triggered) != 1 {
		t.Fatalf("expected 1 triggered alert, got %d", len(triggered))
	}

	// cooldown should prevent immediate re-trigger
	triggered2, _ := eval.EvaluateRules(rules, now)
	if len(triggered2) != 0 {
		t.Fatalf("expected 0 triggered alerts due to cooldown, got %d", len(triggered2))
	}
}


func TestEvaluator_RealWorld(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping real API test in short mode")
	}

	// Real fetcher
	f := fetcher.NewCoinGeckoFetcher()
	eval := NewEvaluator(f)

	// Example alert rules
	rules := []AlertRule{
		{ID: 1, UserID: 42, Coin: "bitcoin", Comparator: Greater, Threshold: 0.01, CooldownMin: 1},
		{ID: 2, UserID: 42, Coin: "ethereum", Comparator: Greater, Threshold: 0.01, CooldownMin: 1},
	}

	now := time.Now()

	triggered, err := eval.EvaluateRules(rules, now)
	if err != nil {
		t.Fatalf("Evaluator failed: %v", err)
	}

	if len(triggered) == 0 {
		t.Fatal("Expected at least one triggered alert")
	}

	for _, a := range triggered {
		t.Logf("Triggered alert: ID=%d, Coin=%s, Threshold=%f", a.ID, a.Coin, a.Threshold)
	}

	// Test cooldown prevents immediate re-trigger
	triggered2, _ := eval.EvaluateRules(rules, now)
	if len(triggered2) != 0 {
		t.Fatalf("Expected 0 triggered alerts due to cooldown, got %d", len(triggered2))
	}
}