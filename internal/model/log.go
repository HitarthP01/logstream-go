package model

//what is this file for?
// This file defines the LogEntry struct, which represents a structured log entry with fields for timestamp, level, and message.
import "time"

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
}
