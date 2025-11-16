// internal/core/format_test.go
package core

import "testing"

func TestFormatAlert(t *testing.T) {
	txt := FormatAlert(AlertRule{ID:1, Coin:"bitcoin", Comparator:Greater, Threshold:1.23})
	if txt == "" { t.Fatal("empty") }
}
