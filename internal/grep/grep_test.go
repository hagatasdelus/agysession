package grep

import (
	"path/filepath"
	"testing"
)

func TestFilter(t *testing.T) {
	t.Setenv("AGYSESSION_CONFIG_DIR", filepath.Join("..", "..", "testdata"))

	// "hello" matches sess-id-1 but not sess-id-2
	hits, err := Filter("hello", Options{Regex: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(hits) != 1 {
		t.Errorf("expected 1 hit, got %d", len(hits))
	}
	if _, ok := hits["sess-id-1"]; !ok {
		t.Errorf("expected sess-id-1 to be matched")
	}

	// Regex matching test
	hitsRegex, err := Filter("run.*", Options{Regex: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hitsRegex) != 1 {
		t.Errorf("expected 1 hit for regex, got %d", len(hitsRegex))
	}
	if _, ok := hitsRegex["sess-id-2"]; !ok {
		t.Errorf("expected sess-id-2 to be matched via regex")
	}

	// Empty query returns nil
	hitsEmpty, err := Filter("", Options{Regex: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hitsEmpty != nil {
		t.Errorf("expected nil for empty query, got %v", hitsEmpty)
	}
}
