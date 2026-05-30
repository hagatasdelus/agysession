package session

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// LabelMaxLen is the max length of the label string.
const LabelMaxLen = 200

// AppDir returns the app directory of Antigravity CLI.
func AppDir() (string, error) {
	if dir := os.Getenv("AGYSESSION_CONFIG_DIR"); dir != "" {
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gemini", "antigravity-cli"), nil
}

// Scan parses the history.jsonl file, dedupes sessions, and returns them
// sorted by last activity (newest first).
func Scan() ([]*Session, error) {
	appDir, err := AppDir()
	if err != nil {
		return nil, err
	}
	historyPath := filepath.Join(appDir, "history.jsonl")
	f, err := os.Open(historyPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	entries := make(map[string]HistoryEntry)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var e HistoryEntry
		if err := json.Unmarshal(sc.Bytes(), &e); err != nil {
			continue
		}
		if e.ConversationID == "" {
			continue
		}
		entries[e.ConversationID] = e
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	var sessions []*Session
	for id, e := range entries {
		// Verify transcript exists
		logPath := filepath.Join(appDir, "brain", id, ".system_generated", "logs", "transcript.jsonl")
		if _, err := os.Stat(logPath); err != nil {
			logPath = filepath.Join(appDir, "brain", id, ".system_generated", "logs", "transcript_full.jsonl")
			if _, err := os.Stat(logPath); err != nil {
				continue
			}
		}

		t := time.UnixMilli(e.Timestamp)
		s := &Session{
			ID:          id,
			CWD:         e.Workspace,
			CWDBasename: filepath.Base(e.Workspace),
			Label:       sanitizeLabel(e.Display),
			LastTime:    t,
			LastEpoch:   t.Unix(),
			CWDExists:   pathIsDir(e.Workspace),
			CWDUnknown:  e.Workspace == "",
			JSONLPath:   logPath,
		}
		if s.CWDBasename == "." || s.CWDBasename == "/" {
			s.CWDBasename = e.Workspace
		}
		sessions = append(sessions, s)
	}

	sort.Slice(sessions, func(i, j int) bool {
		if sessions[i].LastEpoch != sessions[j].LastEpoch {
			return sessions[i].LastEpoch > sessions[j].LastEpoch
		}
		return sessions[i].ID < sessions[j].ID
	})

	return sessions, nil
}

// ScanFiltered behaves like Scan but only includes sessions present in the allow map.
func ScanFiltered(allow map[string]struct{}) ([]*Session, error) {
	sessions, err := Scan()
	if err != nil {
		return nil, err
	}
	// Filter in place
	filtered := sessions[:0]
	for _, s := range sessions {
		if _, ok := allow[s.ID]; ok {
			filtered = append(filtered, s)
		}
	}
	return filtered, nil
}

// FindByID retrieves a session by its conversation ID.
func FindByID(id string) (*Session, error) {
	sessions, err := Scan()
	if err != nil {
		return nil, err
	}
	for _, s := range sessions {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, ErrSessionFileMissing
}

func sanitizeLabel(s string) string {
	s = strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return ' '
		}
		return r
	}, s)
	s = strings.Join(strings.Fields(s), " ")
	return Truncate(s, LabelMaxLen)
}

// Truncate shortens s to at most n runes, appending an ellipsis when cut.
func Truncate(s string, n int) string {
	if n <= 0 {
		return s
	}
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n-1]) + "…"
}

func pathIsDir(path string) bool {
	if path == "" {
		return false
	}
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}
