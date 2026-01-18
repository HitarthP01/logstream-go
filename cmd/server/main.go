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
	err := ingestion.StartTCPIngestion(":9000", proc)
	if err != nil {
		log.Fatal(err)
	}
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
