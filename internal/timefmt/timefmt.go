package timefmt

import (
	"fmt"
	"time"
)

// Parse parses an RFC3339 / RFC3339Nano timestamp string, returning the zero
// time.Time on empty input or any parse error.
func Parse(raw string) time.Time {
	if raw == "" {
		return time.Time{}
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339} {
		if t, err := time.Parse(layout, raw); err == nil {
			return t
		}
	}
	return time.Time{}
}

// Relative formats t relative to now using compact units.
// Falls back to an absolute YYYY-MM-DD date for ages of 7 days or more.
func Relative(t, now time.Time) string {
	if t.IsZero() {
		return "?"
	}
	if t.After(now) {
		t = now
	}
	d := now.Sub(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d/time.Minute))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d/time.Hour))
	case d < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d/(24*time.Hour)))
	default:
		return t.Format("2006-01-02")
	}
}
