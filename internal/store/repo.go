package store

import (
	"context"
	"time"

	"github.com/arenkhachaturian/market-watchdog/internal/core"
)

type AlertRepo interface {
	Create(ctx context.Context, a core.AlertRule) (int, error)
	ListByUser(ctx context.Context, userID int) ([]core.AlertRule, error)
	Delete(ctx context.Context, userID int, alertID int) error

	ListActive(ctx context.Context) ([]core.AlertRule, error)
	UpdateLastNotified(ctx context.Context, alertID int, t time.Time) error
}