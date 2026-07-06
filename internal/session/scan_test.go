package session

import (
	"path/filepath"
	"testing"
)

func TestScan(t *testing.T) {
	t.Setenv("AGYSESSION_CONFIG_DIR", filepath.Join("..", "..", "testdata"))

	sessions, err := Scan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sessions) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(sessions))
	}

	if sessions[0].ID != "sess-id-2" {
		t.Errorf("expected newer session to be sess-id-2, got %s", sessions[0].ID)
	}
	if sessions[0].Label != "Test Prompt 2 Updated" {
		t.Errorf("expected latest label, got %s", sessions[0].Label)
	}
	if sessions[1].ID != "sess-id-1" {
		t.Errorf("expected sess-id-1, got %s", sessions[1].ID)
	}
}

func TestFindByID(t *testing.T) {
	t.Setenv("AGYSESSION_CONFIG_DIR", filepath.Join("..", "..", "testdata"))

	s, err := FindByID("sess-id-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil || s.ID != "sess-id-1" {
		t.Errorf("expected to find sess-id-1, got %+v", s)
	}

	_, err = FindByID("non-existent")
	if err == nil {
		t.Errorf("expected error for non-existent session, got nil")
	}
}

func TestIterContent(t *testing.T) {
	t.Setenv("AGYSESSION_CONFIG_DIR", filepath.Join("..", "..", "testdata"))

	sessions, err := Scan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var messages []string
	err = IterContent(sessions[0].JSONLPath, func(text string) bool {
		messages = append(messages, text)
		return true
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(messages))
	}
	if messages[0] != "running test" || messages[1] != "success" {
		t.Errorf("unexpected messages: %v", messages)
	}
}

func TestCleanAssistantResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Both timestamps",
			input:    "Created At: 2026-06-10T08:02:12Z\nCompleted At: 2026-06-10T08:02:23Z\nHello assistant",
			expected: "Hello assistant",
		},
		{
			name:     "Both timestamps with empty lines",
			input:    "Created At: 2026-06-10T08:02:12Z\nCompleted At: 2026-06-10T08:02:23Z\n\nHello assistant",
			expected: "Hello assistant",
		},
		{
			name:     "No timestamps",
			input:    "Hello assistant\nCreated At: 2026-06-10T08:02:12Z",
			expected: "Hello assistant\nCreated At: 2026-06-10T08:02:12Z",
		},
		{
			name:     "Only timestamps",
			input:    "Created At: 2026-06-10T08:02:12Z\nCompleted At: 2026-06-10T08:02:23Z",
			expected: "",
		},
		{
			name:     "Created At only",
			input:    "Created At: 2026-06-10T08:02:12Z\nHello",
			expected: "Hello",
		},
		{
			name:     "Completed At only",
			input:    "Completed At: 2026-06-10T08:02:23Z\nHello",
			expected: "Hello",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CleanAssistantResponse(tc.input)
			if got != tc.expected {
				t.Errorf("got CleanAssistantResponse(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

