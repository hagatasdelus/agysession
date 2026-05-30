package list

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/hagatasdelus/agysession/internal/ansi"
	"github.com/hagatasdelus/agysession/internal/grep"
	"github.com/hagatasdelus/agysession/internal/session"
	"github.com/hagatasdelus/agysession/internal/timefmt"
)

// Options controls list output.
type Options struct {
	Grep       string
	ExcludeDir string
	Regex      bool
	NoColor    bool
	Color      string
	Out        io.Writer
}

// Run scans sessions, filters them, and writes TSV rows to opts.Out.
func Run(opts Options) error {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}
	sessions, err := loadSessions(opts.Grep, opts.Regex)
	if err != nil {
		return err
	}
	if needle := strings.TrimSpace(opts.ExcludeDir); needle != "" {
		sessions = filterOutByDir(sessions, needle)
	}
	color := colorEnabled(opts)
	now := time.Now()
	w := bufio.NewWriter(opts.Out)
	defer w.Flush()
	for _, s := range sessions {
		if _, err := fmt.Fprintln(w, formatLine(s, now, color)); err != nil {
			return err
		}
	}
	return nil
}

func colorEnabled(opts Options) bool {
	switch strings.ToLower(opts.Color) {
	case "always":
		return true
	case "never":
		return false
	}
	if opts.NoColor {
		return false
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if f, ok := opts.Out.(*os.File); ok {
		return isTerminal(f)
	}
	return false
}

func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func filterOutByDir(sessions []*session.Session, needle string) []*session.Session {
	lneedle := strings.ToLower(needle)
	out := sessions[:0]
	for _, s := range sessions {
		target := s.CWD
		if target == "" {
			target = s.CWDBasename
		}
		if target != "" && strings.Contains(strings.ToLower(target), lneedle) {
			continue
		}
		out = append(out, s)
	}
	return out
}

func loadSessions(query string, regex bool) ([]*session.Session, error) {
	if query == "" {
		return session.Scan()
	}
	allow, err := grep.Filter(query, grep.Options{Regex: regex})
	if err != nil {
		return nil, err
	}
	if len(allow) == 0 {
		return nil, nil
	}
	return session.ScanFiltered(allow)
}

func formatLine(s *session.Session, now time.Time, color bool) string {
	rel := timefmt.Relative(s.LastTime, now)
	base := s.CWDBasename
	if base == "" {
		base = "(unknown)"
	}
	marker := ""
	healthy := true
	switch {
	case s.CWDUnknown:
		marker = "[cwd?] "
		healthy = false
	case !s.CWDExists:
		marker = "[gone] "
		healthy = false
	}
	if color {
		rel = ansi.Dim + padRight(rel, 9) + ansi.Reset
		if healthy {
			base = ansi.Cyan + base + ansi.Reset
		} else {
			base = ansi.Yellow + base + ansi.Reset
			marker = ansi.Yellow + marker + ansi.Reset
		}
	} else {
		rel = padRight(rel, 9)
	}
	return fmt.Sprintf("%s\t%d\t%s\t%s\t%s%s",
		s.ID, s.LastEpoch, rel, base, marker, s.Label)
}

func padRight(s string, n int) string {
	for len([]rune(s)) < n {
		s += " "
	}
	return s
}
