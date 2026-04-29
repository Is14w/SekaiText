package service

import (
	"sync"
	"time"
)

// LogEntry represents a single server log entry.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

// LogBuffer is a thread-safe ring buffer for server logs.
type LogBuffer struct {
	mu    sync.Mutex
	lines []LogEntry
	cap   int
}

// NewLogBuffer creates a log buffer with the given capacity.
func NewLogBuffer(capacity int) *LogBuffer {
	return &LogBuffer{
		lines: make([]LogEntry, 0, capacity),
		cap:   capacity,
	}
}

// Write appends a log entry.
func (lb *LogBuffer) Write(msg string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	if len(lb.lines) >= lb.cap {
		lb.lines = lb.lines[1:]
	}
	lb.lines = append(lb.lines, LogEntry{
		Timestamp: time.Now().Format("15:04:05"),
		Message:   msg,
	})
}

// Lines returns all buffered log entries.
func (lb *LogBuffer) Lines() []LogEntry {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	out := make([]LogEntry, len(lb.lines))
	copy(out, lb.lines)
	return out
}
