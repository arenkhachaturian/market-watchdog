package core

import "time"

type Comparator int

const (
	Greater Comparator = iota
	Less
	PercentChange
)

type AlertRule struct {
	ID 				int
	UserID  		int
	Coin			string
	Comparator		Comparator
	Threshold		float64
	LastTriggered 	*time.Time
	CooldownMin 	int
}

func (a *AlertRule) CanTrigger(now time.Time) bool {
	if a.LastTriggered == nil {
		return true
	}
	elapsed := now.Sub(*a.LastTriggered)
	return elapsed.Minutes() >= float64(a.CooldownMin)
}
