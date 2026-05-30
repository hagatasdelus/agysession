package list

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	t.Setenv("AGYSESSION_CONFIG_DIR", filepath.Join("..", "..", "testdata"))

	var buf bytes.Buffer
	err := Run(Options{
		Grep:       "",
		ExcludeDir: "",
		Regex:      false,
		Color:      "never",
		NoColor:    true,
		Out:        &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines of TSV output, got %d", len(lines))
	}

	// 1st line (newest session should be sess-id-2)
	if !strings.HasPrefix(lines[0], "sess-id-2") {
		t.Errorf("expected newest session sess-id-2 first, got: %s", lines[0])
	}
	if !strings.Contains(lines[0], "project2") {
		t.Errorf("expected path project2 in output, got: %s", lines[0])
	}

	// Test ExcludeDir
	buf.Reset()
	err = Run(Options{
		Grep:       "",
		ExcludeDir: "project1",
		Regex:      false,
		Color:      "never",
		NoColor:    true,
		Out:        &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output = buf.String()
	lines = strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 1 || !strings.HasPrefix(lines[0], "sess-id-2") {
		t.Errorf("expected only sess-id-2 after excluding project1, got: %s", output)
	}

	// Test Grep filter in List
	buf.Reset()
	err = Run(Options{
		Grep:       "hello",
		ExcludeDir: "",
		Regex:      false,
		Color:      "never",
		NoColor:    true,
		Out:        &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output = buf.String()
	lines = strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 1 || !strings.HasPrefix(lines[0], "sess-id-1") {
		t.Errorf("expected only sess-id-1 matching 'hello', got: %s", output)
	}
}
