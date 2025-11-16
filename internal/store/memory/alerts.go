package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/arenkhachaturian/market-watchdog/internal/core"
	"github.com/arenkhachaturian/market-watchdog/internal/store"
)

type Alerts struct {
	mu 				   sync.Mutex
	seq 			   int
	all 			   []core.AlertRule
	defaultCooldownMin int
}

func NewAlerts(defaultCooldownMin int) *Alerts {
	return &Alerts{defaultCooldownMin: defaultCooldownMin}
}

func (repo *Alerts) Create(_ context.Context, alert core.AlertRule) (int, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.seq++
	alert.ID = repo.seq
	if alert.CooldownMin == 0 {
		alert.CooldownMin = repo.defaultCooldownMin
	}
	repo.all = append(repo.all, alert)
	return alert.ID, nil
}

func (repo *Alerts) ListByUser(_ context.Context, userID int) ([]core.AlertRule, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	out := make([]core.AlertRule, 0, len(repo.all))
	for _, alert := range repo.all {
		if alert.UserID == userID {
			out = append(out, alert)
		}
	}
	return out, nil
}

func (repo *Alerts) Delete(_ context.Context, userID int, alertID int) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	result := make([]core.AlertRule, 0, len(repo.all))
	found := false
	for _, alert := range repo.all {
		if alert.UserID == userID && alert.ID == alertID {
			found = true
			continue
		}
		result = append(result, alert)
	}
	if found != true {
		return errors.New("not found")
	}
	repo.all = result
	return nil
}

// refactor after move on go 1.21
func (r *Alerts) ListActive(_ context.Context) ([]core.AlertRule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]core.AlertRule, len(r.all))
	copy(out, r.all)
	return out, nil
}

func (repo *Alerts) UpdateLastNotified(_ context.Context, alertID int, t time.Time) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for i, _ := range repo.all {
		if repo.all[i].ID == alertID {
			repo.all[i].LastTriggered = &t
			return nil
		}
	}
	return errors.New("not found")
}

var _ store.AlertRepo = (*Alerts)(nil);