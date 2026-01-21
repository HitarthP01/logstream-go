package processor

import (
	"go-sentinel/internal/model"
	// "sync"
)

// Processor defines the interface for processing log entries

type Processor struct {
	
	Input chan model.LogEntry // channel to receive log entries
	//explain this syntax
	// Input is a channel of type model.LogEntry, which means it can send and receive log entries.
	// It is used to receive log entries from the ingestion layer and pass them to the processor for processing.
	// Channels are a way to communicate between different parts of a program.
	// In this case, the ingestion layer sends log entries to the processor through this channel.
	// The channel is unbuffered, meaning it can only hold one log entry at a time.
	// When the channel is full, the ingestion layer will block until there is space in the channel.
	// This ensures that the processor can keep up with the ingestion rate.
	
	// mu            sync.Mutex // to ensure thread-safe access to shared resources
	Logs          []model.LogEntry
	LevelCounts   map[string]int
	SeverityCount map[string]int
}

//have no idea what this func is doing and what is this syntax
// This function named New initializes and returns a pointer to a Processor struct.
// It sets up the LevelCounts and SeverityCount maps to be empty maps ready for use.

// initiator func
func New() *Processor {
	return &Processor{
		Input:         make(chan model.LogEntry,100),
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

func (p *Processor) Run() {
	for entry := range p.Input {
		p.Process(entry)
	}
}
func (p *Processor) Process(entry model.LogEntry) {
	// p.mu.Lock()
	// defer p.mu.Unlock()
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
