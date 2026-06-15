package preview

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hagatasdelus/agysession/internal/ansi"
	"github.com/hagatasdelus/agysession/internal/session"
	"github.com/hagatasdelus/agysession/internal/timefmt"
)

// Options controls preview rendering.
type Options struct {
	Query string
	Regex bool
}

const (
	maxMessages = 30
	maxBodyLen  = 200
)

type messageItem struct {
	Role      string
	Timestamp time.Time
	Body      string
}

// Run writes the preview pane content for a given session id to stdout.
func Run(id string, opts Options) error {
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
	return render(s, os.Stdout, opts)
}

func render(s *session.Session, out io.Writer, opts Options) error {
	messages, startedAt, totalMsgs, err := loadMessages(s.JSONLPath)
	if err != nil {
		return err
	}
	if startedAt.IsZero() && !s.LastTime.IsZero() {
		startedAt = s.LastTime
	}

	now := time.Now()
	w := bufio.NewWriter(out)
	defer w.Flush()

	fmt.Fprintf(w, "%ssession%s : %s\n", ansi.Bold, ansi.Reset, s.ID)
	cwd := s.CWD
	if cwd == "" {
		cwd = "(unknown)"
	}
	if !s.CWDExists {
		cwd = ansi.Yellow + cwd + " [gone]" + ansi.Reset
	} else {
		cwd = ansi.Cyan + cwd + ansi.Reset
	}
	fmt.Fprintf(w, "%sproject%s : %s\n", ansi.Bold, ansi.Reset, cwd)
	fmt.Fprintf(w, "%sstarted%s : %s  %s(%s)%s\n",
		ansi.Bold, ansi.Reset,
		startedAt.Local().Format("2006-01-02 15:04"),
		ansi.Dim, relativeOrFuture(startedAt, now), ansi.Reset,
	)
	fmt.Fprintf(w, "%slast%s    : %s  %s(%d msgs)%s\n",
		ansi.Bold, ansi.Reset,
		s.LastTime.Local().Format("2006-01-02 15:04"),
		ansi.Dim, totalMsgs, ansi.Reset,
	)
	fmt.Fprintln(w, ansi.Dim+strings.Repeat("─", 60)+ansi.Reset)

	tail := messages
	if len(tail) > maxMessages {
		tail = tail[len(tail)-maxMessages:]
	}
	for _, m := range tail {
		writeMessage(w, m, opts)
	}
	return nil
}

func writeMessage(w io.Writer, m messageItem, opts Options) {
	role := m.Role
	color := ansi.Green
	if role == "assistant" {
		role = "asst"
		color = ansi.Cyan
	}
	stamp := "--:--"
	if !m.Timestamp.IsZero() {
		stamp = m.Timestamp.Local().Format("15:04")
	}
	body := highlightMatches(truncateBody(m.Body), opts)
	fmt.Fprintf(w, "%s[%s %s]%s %s\n", color, role, stamp, ansi.Reset, body)
}

func highlightMatches(s string, opts Options) string {
	if strings.TrimSpace(opts.Query) == "" {
		return s
	}
	pattern := opts.Query
	if !opts.Regex {
		pattern = regexp.QuoteMeta(opts.Query)
	}
	re, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		return s
	}
	locs := re.FindAllStringIndex(s, -1)
	if len(locs) == 0 {
		return s
	}
	var b strings.Builder
	last := 0
	for _, loc := range locs {
		if loc[0] == loc[1] {
			continue
		}
		b.WriteString(s[last:loc[0]])
		b.WriteString(ansi.Highlight)
		b.WriteString(s[loc[0]:loc[1]])
		b.WriteString(ansi.Reset)
		last = loc[1]
	}
	b.WriteString(s[last:])
	return b.String()
}

func truncateBody(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	lines := strings.Split(s, "\n")
	var kept []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		kept = append(kept, l)
		if len(kept) >= 2 {
			break
		}
	}
	joined := strings.Join(kept, " | ")
	return session.Truncate(joined, maxBodyLen)
}

const previewLineCap = 16 * 1024 * 1024

func loadMessages(path string) ([]messageItem, time.Time, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, time.Time{}, 0, err
	}
	defer f.Close()

	r := bufio.NewReaderSize(f, 64*1024)
	var (
		ring      = make([]messageItem, maxMessages)
		startedAt time.Time
		total     int
	)
	for {
		line, err := readJSONLLine(r, previewLineCap)
		if line != "" {
			if item, ts, ok := parseMessageLine(line); ok {
				if startedAt.IsZero() && !ts.IsZero() {
					startedAt = ts
				}
				ring[total%maxMessages] = item
				total++
			}
		}
		if errors.Is(err, io.EOF) {
			return collectRing(ring, total), startedAt, total, nil
		}
		if err != nil {
			return collectRing(ring, total), startedAt, total, err
		}
	}
}

func collectRing(ring []messageItem, total int) []messageItem {
	n := min(total, maxMessages)
	if n == 0 {
		return nil
	}
	out := make([]messageItem, 0, n)
	start := 0
	if total > maxMessages {
		start = total % maxMessages
	}
	for i := range n {
		out = append(out, ring[(start+i)%maxMessages])
	}
	return out
}

func readJSONLLine(r *bufio.Reader, max int) (string, error) {
	var (
		buf       strings.Builder
		truncated bool
	)
	for {
		chunk, err := r.ReadSlice('\n')
		if len(chunk) > 0 {
			if !truncated {
				if buf.Len()+len(chunk) > max {
					truncated = true
				} else {
					buf.Write(chunk)
				}
			}
		}
		if errors.Is(err, bufio.ErrBufferFull) {
			truncated = true
			continue
		}
		if truncated {
			return "", err
		}
		return strings.TrimSpace(buf.String()), err
	}
}

func parseMessageLine(line string) (messageItem, time.Time, bool) {
	var e session.TranscriptEntry
	if err := json.Unmarshal([]byte(line), &e); err != nil {
		return messageItem{}, time.Time{}, false
	}
	if e.Source != "USER_EXPLICIT" && e.Source != "MODEL" {
		return messageItem{}, time.Time{}, false
	}
	if e.Content == "" {
		return messageItem{}, time.Time{}, false
	}
	role := "user"
	body := e.Content
	switch e.Source {
	case "USER_EXPLICIT":
		body = session.CleanUserRequest(body)
	case "MODEL":
		role = "assistant"
		body = session.CleanAssistantResponse(body)
	}
	ts := timefmt.Parse(e.CreatedAt)
	return messageItem{Role: role, Timestamp: ts, Body: body}, ts, true
}


func relativeOrFuture(t, now time.Time) string {
	if !t.IsZero() && t.After(now) {
		return "in the future"
	}
	return timefmt.Relative(t, now)
}
