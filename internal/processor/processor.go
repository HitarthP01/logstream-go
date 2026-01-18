package processor

import (
	"go-sentinel/internal/model"
)

// Processor defines the interface for processing log entries

type Processor struct {
	Logs          []model.LogEntry
	LevelCounts   map[string]int
	SeverityCount map[string]int
}

//have no idea what this func is doing and what is this syntax
// This function named New initializes and returns a pointer to a Processor struct.
// It sets up the LevelCounts and SeverityCount maps to be empty maps ready for use.

func New() *Processor {
	return &Processor{
		LevelCounts:   make(map[string]int),
		SeverityCount: make(map[string]int),
	}
}

//explain the syntax of this func
// This method named Process is defined on the Processor struct.
// It takes a LogEntry as a parameter and processes it.
// The actual processing logic is not implemented yet, as the function body is currently empty.

// (p *Processor) indicates that this is a method with a receiver of type *Processor,
//
//	allowing it to access and modify the Processor's fields.
func (p *Processor) Process(entry model.LogEntry) {
	p.Logs = append(p.Logs, entry)

	// Update level counts
	p.LevelCounts[entry.Level]++
	// Update severity counts based on log level

	SeverityCount := classifySeverity(entry.Level)
	p.SeverityCount[SeverityCount]++

}

func classifySeverity(level string) string {
	switch level {
	case "ERROR":
		return "CRITICAL"
	case "WARN":
		return "WARNING"
	default:
		return "INFO"
	}
}
