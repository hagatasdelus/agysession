package preview

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hagatasdelus/agysession/internal/session"
)

func TestRender(t *testing.T) {
	t.Setenv("AGYSESSION_CONFIG_DIR", filepath.Join("..", "..", "testdata"))

	s, err := session.FindByID("sess-id-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	err = render(s, &buf, Options{Query: "hello", Regex: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, s.ID) {
		t.Errorf("expected output to contain session ID %s, got: %s", s.ID, output)
	}
	if !strings.Contains(output, "/mock/project1") {
		t.Errorf("expected project path, got: %s", output)
	}
	// "hello" should be highlighted (reverse video ANSI escape: \x1b[7m ... \x1b[0m)
	if !strings.Contains(output, "\x1b[7mhello\x1b[0m assistant") {
		t.Errorf("expected highlighted query, got: %q", output)
	}
}
