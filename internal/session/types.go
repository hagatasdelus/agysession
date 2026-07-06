package session

import (
	"errors"
	"time"
)

// ErrSessionFileMissing is returned when a session transcript file doesn't exist.
var ErrSessionFileMissing = errors.New("session file missing")

// ErrSessionEmpty is returned when the session has no content.
var ErrSessionEmpty = errors.New("session has no usable content")

// Session is a parsed view of one Antigravity CLI session.
type Session struct {
	ID          string
	CWD         string
	CWDBasename string
	Label       string
	LastTime    time.Time
	LastEpoch   int64
	CWDExists   bool
	CWDUnknown  bool
	JSONLPath   string
}

// HistoryEntry matches a row in ~/.gemini/antigravity-cli/history.jsonl.
type HistoryEntry struct {
	Display        string `json:"display"`
	Timestamp      int64  `json:"timestamp"` // epoch ms
	Workspace      string `json:"workspace"`
	ConversationID string `json:"conversationId"`
}

// TranscriptEntry matches a row in transcript.jsonl.
type TranscriptEntry struct {
	StepIndex int    `json:"step_index"`
	Source    string `json:"source"` // MODEL, USER_EXPLICIT, SYSTEM
	Type      string `json:"type"`   // USER_INPUT, PLANNER_RESPONSE, etc.
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"` // RFC3339 timestamp
	Content   string `json:"content"`
}

// IsUserRequest returns true if the entry represents a user request.
func (e *TranscriptEntry) IsUserRequest() bool {
	return e.Source == "USER_EXPLICIT" && e.Type == "USER_INPUT"
}

// IsAssistantResponse returns true if the entry represents an assistant response.
func (e *TranscriptEntry) IsAssistantResponse() bool {
	return e.Source == "MODEL" && e.Type == "PLANNER_RESPONSE"
}
