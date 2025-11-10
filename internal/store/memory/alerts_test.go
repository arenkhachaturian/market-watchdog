// internal/store/memory/alerts_test.go
package memory

import (
	"context"
	"testing"
	"time"

	"github.com/arenkhachaturian/market-watchdog/internal/core"
)

func TestAlertsRepo(t *testing.T) {
	alertRepo := NewAlerts()
	ctx := context.Background()

	// create two alerts
	t.Logf("start")
	firstID, err := alertRepo.Create(ctx, core.AlertRule{
		UserID: 7, Coin: "bitcoin", Comparator: core.Greater, Threshold: 1, CooldownMin: 30,
	})
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	t.Logf("created alert id=%d (user=%d coin=%s)", firstID, 7, "bitcoin")

	secondID, err := alertRepo.Create(ctx, core.AlertRule{
		UserID: 7, Coin: "eth", Comparator: core.Less, Threshold: 2, CooldownMin: 15,
	})
	if err != nil {
		t.Fatalf("create second: %v", err)
	}
	t.Logf("created alert id=%d (user=%d coin=%s)", secondID, 7, "eth")

	// list by user
	userAlerts, err := alertRepo.ListByUser(ctx, 7)
	if err != nil {
		t.Fatalf("list by user: %v", err)
	}
	t.Logf("list by user -> count=%d", len(userAlerts))
	if len(userAlerts) != 2 {
		t.Fatalf("want 2, got %d", len(userAlerts))
	}

	// delete second
	if err := alertRepo.Delete(ctx, 7, secondID); err != nil {
		t.Fatalf("delete id=%d: %v", secondID, err)
	}
	t.Logf("deleted alert id=%d", secondID)

	// list again
	userAlerts, _ = alertRepo.ListByUser(ctx, 7)
	t.Logf("list after delete -> count=%d", len(userAlerts))
	if len(userAlerts) != 1 {
		t.Fatalf("want 1 after delete, got %d", len(userAlerts))
	}

	// update last-notified
	now := time.Now()
	if err := alertRepo.UpdateLastNotified(ctx, firstID, now); err != nil {
		t.Fatalf("update last-notified id=%d: %v", firstID, err)
	}
	allActive, _ := alertRepo.ListActive(ctx)
	t.Logf("list active -> count=%d", len(allActive))
	if allActive[0].LastTriggered == nil {
		t.Fatal("LastTriggered not updated")
	}
	t.Logf("alert id=%d lastTriggered=%s", allActive[0].ID, allActive[0].LastTriggered.UTC().Format(time.RFC3339))
}
