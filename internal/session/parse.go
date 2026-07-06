package session

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strings"
)

// IterContent reads an Antigravity CLI transcript.jsonl and invokes fn with the
// content of each user/assistant step. fn can return false to stop early.
func IterContent(path string, fn func(text string) bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReaderSize(f, 64*1024)
	for {
		line, err := readLine(r)
		if line != "" {
			var e TranscriptEntry
			if err := json.Unmarshal([]byte(line), &e); err == nil {
				if (e.IsUserRequest() || e.IsAssistantResponse()) && e.Content != "" {
					content := e.Content
					switch e.Source {
					case "USER_EXPLICIT":
						content = CleanUserRequest(content)
					case "MODEL":
						content = CleanAssistantResponse(content)
					}
					if !fn(content) {
						return nil
					}
				}
			}
		}
		if err != nil {
			break
		}
	}
	return nil
}

func readLine(r *bufio.Reader) (string, error) {
	var buf strings.Builder
	for {
		chunk, err := r.ReadSlice('\n')
		if len(chunk) > 0 {
			buf.Write(chunk)
		}
		if errors.Is(err, bufio.ErrBufferFull) {
			continue
		}
		return strings.TrimSpace(buf.String()), err
	}
}

// CleanUserRequest extracts the core request text from a USER_EXPLICIT message content.
// It removes <USER_REQUEST> ... </USER_REQUEST> tags and metadata if present.
func CleanUserRequest(s string) string {
	const startTag = "<USER_REQUEST>"
	const endTag = "</USER_REQUEST>"

	startIdx := strings.Index(s, startTag)
	if startIdx == -1 {
		return s
	}

	contentStart := startIdx + len(startTag)
	endIdx := strings.Index(s[contentStart:], endTag)
	if endIdx == -1 {
		return strings.TrimSpace(s[contentStart:])
	}

	return strings.TrimSpace(s[contentStart : contentStart+endIdx])
}

// CleanAssistantResponse removes metadata lines (like Created At/Completed At) from assistant messages.
func CleanAssistantResponse(s string) string {
	lines := strings.Split(s, "\n")
	var startIdx int
	for startIdx < len(lines) {
		trimmed := strings.TrimSpace(lines[startIdx])
		if trimmed == "" {
			startIdx++
			continue
		}
		if strings.HasPrefix(trimmed, "Created At:") || strings.HasPrefix(trimmed, "Completed At:") {
			startIdx++
			continue
		}
		break
	}
	return strings.Join(lines[startIdx:], "\n")
}


