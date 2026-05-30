package resume

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/hagatasdelus/agysession/internal/session"
)

var execve = syscall.Exec

// Run changes the directory to the original session CWD and execs "agy --conversation <id>"
// replacing the current process.
func Run(id string) error {
	s, err := session.FindByID(id)
	if err != nil {
		if errors.Is(err, session.ErrSessionFileMissing) {
			return fmt.Errorf("session not found: %s", id)
		}
		return err
	}
	if s == nil {
		return fmt.Errorf("session not found: %s", id)
	}
	if s.CWDUnknown || s.CWD == "" {
		return fmt.Errorf("cwd unknown for session %s; refusing to resume", id)
	}
	if !s.CWDExists {
		return fmt.Errorf("original cwd is gone: %s", s.CWD)
	}

	agyPath, err := exec.LookPath("agy")
	if err != nil {
		// Fallback to ~/.local/bin/agy
		home, err2 := os.UserHomeDir()
		if err2 == nil {
			fallback := filepath.Join(home, ".local", "bin", "agy")
			if _, err3 := os.Stat(fallback); err3 == nil {
				agyPath = fallback
			}
		}
		if agyPath == "" {
			return fmt.Errorf("agy CLI not found in PATH: %w", err)
		}
	}

	if err := os.Chdir(s.CWD); err != nil {
		return fmt.Errorf("chdir to %s: %w", s.CWD, err)
	}

	return execve(agyPath, []string{"agy", "--conversation", id}, os.Environ())
}
