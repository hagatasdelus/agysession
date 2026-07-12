package resume

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("AGYSESSION_CONFIG_DIR", tempDir)

	// Create test index pointing to a real workspace dir inside tempDir
	workspaceDir := filepath.Join(tempDir, "myproject")
	if err := os.Mkdir(workspaceDir, 0755); err != nil {
		t.Fatalf("failed to create mock workspace: %v", err)
	}

	historyContent := fmt.Sprintf(
		`{"display":"Test","timestamp":1780000000000,"workspace":%q,"conversationId":"sess-resume-1"}`+"\n",
		workspaceDir,
	)
	if err := os.WriteFile(filepath.Join(tempDir, "history.jsonl"), []byte(historyContent), 0600); err != nil {
		t.Fatalf("failed to write history.jsonl: %v", err)
	}

	// Create dummy transcript
	logDir := filepath.Join(tempDir, "brain", "sess-resume-1", ".system_generated", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("failed to create log dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(logDir, "transcript.jsonl"), []byte("{}\n"), 0600); err != nil {
		t.Fatalf("failed to write transcript: %v", err)
	}

	var called int
	var calledArgv0 string
	var calledArgv []string
	originalExecve := execve
	originalLookPath := lookPath
	defer func() {
		execve = originalExecve
		lookPath = originalLookPath
	}()

	execve = func(argv0 string, argv []string, envv []string) error {
		called++
		calledArgv0 = argv0
		calledArgv = argv
		return nil
	}
	lookPath = func(file string) (string, error) {
		return "/mock/bin/agy", nil
	}

	err := Run("sess-resume-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if called != 1 {
		t.Errorf("expected 1 call, got %d", called)
	}
	if calledArgv0 == "" {
		t.Errorf("expected non-empty executable path")
	}
	if len(calledArgv) < 3 || calledArgv[1] != "--conversation" || calledArgv[2] != "sess-resume-1" {
		t.Errorf("unexpected exec arguments: %v", calledArgv)
	}
}
