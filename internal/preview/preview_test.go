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

func TestParseMessageLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantRole string
		wantBody string
		wantOk   bool
	}{
		{
			name:     "User message",
			line:     `{"step_index":0,"source":"USER_EXPLICIT","type":"USER_INPUT","status":"DONE","created_at":"2026-05-30T00:00:00Z","content":"<USER_REQUEST>\nhello"}`,
			wantRole: "user",
			wantBody: "hello",
			wantOk:   true,
		},
		{
			name:     "Assistant message with timestamps",
			line:     `{"step_index":1,"source":"MODEL","type":"PLANNER_RESPONSE","status":"DONE","created_at":"2026-05-30T00:01:00Z","content":"Created At: 2026-06-10T08:02:12Z\nCompleted At: 2026-06-10T08:02:23Z\nHello assistant"}`,
			wantRole: "assistant",
			wantBody: "Hello assistant",
			wantOk:   true,
		},
		{
			name:     "Assistant message from command execution (should be ignored)",
			line:     `{"step_index":1,"source":"MODEL","type":"RUN_COMMAND","status":"DONE","created_at":"2026-05-30T00:02:30Z","content":"The command completed successfully."}`,
			wantOk:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			item, _, ok := parseMessageLine(tc.line)
			if ok != tc.wantOk {
				t.Fatalf("expected ok=%v, got %v", tc.wantOk, ok)
			}
			if !ok {
				return
			}
			if item.Role != tc.wantRole {
				t.Errorf("expected role %q, got %q", tc.wantRole, item.Role)
			}
			if item.Body != tc.wantBody {
				t.Errorf("expected body %q, got %q", tc.wantBody, item.Body)
			}
		})
	}
}

