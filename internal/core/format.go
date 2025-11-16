package core

import "fmt"

func FormatAlert(a AlertRule) string {
	return fmt.Sprintf("Alert #%d: %s %v %.4f",
		a.ID, a.Coin, a.Comparator, a.Threshold)
}
