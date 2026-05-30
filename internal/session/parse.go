package session

import (
	"bufio"
	"encoding/json"
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
				if (e.Source == "USER_EXPLICIT" || e.Source == "MODEL") && e.Content != "" {
					if !fn(e.Content) {
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
		if err == bufio.ErrBufferFull {
			continue
		}
		return strings.TrimSpace(buf.String()), err
	}
}
