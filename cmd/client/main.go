package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:9000")

	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println()

	logs := []string{
		"2026-01-08 10:23:45 INFO User login successful user_id=1234",
		"2026-01-08 10:24:10 ERROR Database connection failed",
		"2026-01-08 10:25:05 WARN Disk space running low user_id=1234",
		"2026-01-08 10:26:30 INFO File uploaded successfully file_id=5678",
		"2026-01-08 10:27:15 ERROR Timeout while processing request",
	}

	for _, log := range logs { //what does _, log := range logs mean?
		// In the for loop, 'range logs' iterates over each element in the 'logs' slice.
		// The underscore '_' is used to ignore the index of the current element, since we don't need it.
		// 'log' is the variable that holds the value of the current element in each iteration.

		fmt.Fprintln(conn, log) // what is Fprintln?
		// fmt.Fprintln is a function from the fmt package that writes formatted output to a specified writer (in this case, the TCP connection 'conn').
		// It appends a newline character at the end of the string, making it suitable for sending log entries line by line.

		//how is it different from Println?
		// The main difference between fmt.Fprintln and fmt.Println is that fmt.Fprintln allows you to specify the destination (writer) for the output,
		// while fmt.Println always writes to the standard output (console).

		/*
			┌────────┐    fmt.Fprintln(conn, log)    ┌────────┐
			│ Client │ ─────────────────────────────▶│ Server │
			│        │         TCP Connection        │ :9000  │
			└────────┘                               └────────┘

		*/
		time.Sleep(1 * time.Second) // Simulate delay between log entries
	}
}

/*

1. Client starts
         │
         ▼
2. Connect to server (net.Dial)
         │
         ▼
3. Schedule cleanup (defer conn.Close)
         │
         ▼
4. Loop through logs ◄────────┐
         │                    │
         ▼                    │
5. Send log to server         │
         │                    │
         ▼                    │
6. Wait 1 second              │
         │                    │
         ▼                    │
7. More logs? ────── Yes ─────┘
         │
         No
         ▼
8. Function exits → defer runs → connection closes
*/
