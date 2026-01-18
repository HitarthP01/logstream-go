package parser

import (
	"go-sentinel/internal/model"
	// "internal/model"
	"strings"
	"time"
)

// what is this syntax of go func?
// This syntax defines a function named ParseLogLine that takes a string parameter 'line' and returns a LogEntry struct and a boolean value.

// this func returns multiple values: logentry and bool
// sytax: func fnname(parameter type) (ret val1, ret val2){}
func ParseLogLine(line string) (model.LogEntry, bool) {
	parts := strings.SplitN(line, " ", 4)
	if len(parts) < 4 {
		return model.LogEntry{}, false
	}

	timestr := parts[0] + " " + parts[1]
	timestamp, err := time.Parse("2006-01-02 15:04:05", timestr)
	if err != nil {
		return model.LogEntry{}, false
	}

	return model.LogEntry{
		Timestamp: timestamp,
		Level:     parts[2],
		Message:   parts[3],
	}, true

}
