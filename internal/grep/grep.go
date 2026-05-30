package grep

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/hagatasdelus/agysession/internal/session"
)

// Options controls Filter behavior.
type Options struct {
	Regex bool
}

// Filter scans session transcripts and returns the set of conversation IDs
// whose user/model content matches the query.
func Filter(query string, opts Options) (map[string]struct{}, error) {
	if strings.TrimSpace(query) == "" {
		return nil, nil
	}

	match, err := buildMatcher(query, opts)
	if err != nil {
		return nil, err
	}

	sessions, err := session.Scan()
	if err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return map[string]struct{}{}, nil
	}

	set := make(map[string]struct{})
	var mu sync.Mutex
	concurrency := max(runtime.NumCPU()*2, 4)
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for _, s := range sessions {
		wg.Add(1)
		sem <- struct{}{}
		go func(sess *session.Session) {
			defer wg.Done()
			defer func() { <-sem }()
			var hit bool
			_ = session.IterContent(sess.JSONLPath, func(text string) bool {
				if match(text) {
					hit = true
					return false
				}
				return true
			})
			if hit {
				mu.Lock()
				set[sess.ID] = struct{}{}
				mu.Unlock()
			}
		}(s)
	}
	wg.Wait()
	return set, nil
}

func buildMatcher(query string, opts Options) (func(string) bool, error) {
	if opts.Regex {
		re, err := regexp.Compile("(?i)" + query)
		if err != nil {
			return nil, fmt.Errorf("invalid regex: %w", err)
		}
		return re.MatchString, nil
	}
	needle := strings.ToLower(query)
	return func(text string) bool {
		return strings.Contains(strings.ToLower(text), needle)
	}, nil
}
