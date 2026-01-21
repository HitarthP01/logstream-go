package main

import (
	"go-sentinel/internal/ingestion"
	"go-sentinel/internal/processor"
	"log"
)

func main() {

	//why do we have to create a new processor here?
	// We create a new processor here to handle the processing of log entries
	// that will be ingested by the TCP server. The processor is responsible for
	// storing the logs and maintaining counts of log levels and severities.
	// By creating a new processor instance, we ensure that we have a fresh state
	// to work with when the server starts receiving log data.
	proc := processor.New()
  go proc.Run()
	err := ingestion.StartTCPIngestion(":9000", proc)
	if err != nil {
		log.Fatal(err)
	}

	//how many processors will be created if multiple clients connect to the server?
	// Only one processor instance is created in the main function.
	// This single processor instance is shared among all client connections.
	// Each time a new client connects, a new goroutine is spawned to handle that connection,
	// but they all use the same processor instance passed to the StartTCPIngestion function.

	//so, all log entries from different clients will be processed by the same processor,and share resources like Logs slice and LevelCounts map?
	// Yes, all log entries from different clients will be processed by the same processor instance.
	// This means that they will share resources like the Logs slice and LevelCounts map.
	// As a result, the processor will aggregate log data from all connected clients,
	// allowing for a consolidated view of log entries and their statistics.
}

// 			┌────────────┐                      ┌────────────┐
// 			│   Client   │                      │   Server   │
// 			│            │                      │            │
// 			│ net.Dial() │ ──── connects ────▶  │  Accept()  │
// 			│     │      │                      │     │      │
// 			│     ▼      │                      │     ▼      │
// 			│   conn     │ ◄═══════════════════▶│   conn     │
// 			└────────────┘    TCP Connection    └────────────┘

/*
        CLIENT                                    SERVER
        ══════                                    ══════
           │                                         │
           │                            1. net.Listen(":9000")
           │                                         │
           │                            2. listener.Accept() ◄── BLOCKS
           │                                         │
3. net.Dial("localhost:9000") ─────── connects ─────▶│
           │                                         │
           │◄═══════ TCP CONNECTION ESTABLISHED ═════▶│
           │                                         │
           │                            4. bufio.NewScanner(conn)
           │                                         │
           │                            5. scanner.Scan() ◄── BLOCKS
           │                                         │
6. fmt.Fprintln(conn, log1) ──────── sends data ────▶│
           │                                         │
           │                            7. Receives log1, prints it
           │                                         │
           │                            8. scanner.Scan() ◄── BLOCKS
           │                                         │
9. time.Sleep(1 second)                              │ (waiting)
           │                                         │
10. fmt.Fprintln(conn, log2) ─────── sends data ────▶│
           │                                         │
           │                            11. Receives log2, prints it
           │                                         │
         ...                                       ...
           │                                         │
12. Loop ends                                        │
           │                                         │
13. defer conn.Close() ──────── closes connection ──▶│
           │                                         │
           │                            14. scanner.Scan() returns false
           │                                         │
           │                            15. Loop ends, prints "Connection closed"


*/

/*

Time ──────────────────────────────────────────────────────────────▶

SERVER:  [Listen] [Accept...waiting...] [Scan...waiting...] [Print] [Scan...waiting...] [Print] [Done]
                         │                      │                          │
                         │                      │                          │
CLIENT:            [Dial]│              [Send]  │  [Sleep]  [Send]         │  [Close]
                         │                      │                          │
                    Connection              Data flows              Connection
                    Established                                      Closed
*/
