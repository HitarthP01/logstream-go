package ingestion

import (
	"bufio"
	"fmt"
	"go-sentinel/internal/parser"
	"go-sentinel/internal/processor"
	"net"
	// "text/template/parse"
)

// StartTCPIngestion starts a TCP server to ingest log data
// return type is error, because starting a server may fail due to various reasons (e.g., port already in use).
func StartTCPIngestion(address string, proc *processor.Processor) error {
	listenr, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listenr.Close()
	fmt.Println("TCP server started on", address)
	for {
		conn, err := listenr.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		//maintain concurrency with goroutines
		go handleConnection(conn, proc)
	}

}

func handleConnection(conn net.Conn, proc *processor.Processor) {

	//why do we close the conn here?
	// We use defer to ensure that the connection is closed when the function exits,
	// regardless of whether it exits normally or due to an error.
	// This helps to prevent resource leaks by making sure that the connection is properly closed after we're done using it.

	//why to use defer? if we can just close it at the end of the func?
	// Using defer ensures that the connection will be closed even if an error occurs or if the function returns early.
	// If we were to close it only at the end of the function, we might miss closing it in case of an error or early return,
	// leading to potential resource leaks.
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()

		parsedEntry, ok := parser.ParseLogLine(line)
		if !ok {
			fmt.Println("Failed to parse log line:", line)
			continue
		}
		proc.Process(parsedEntry)
		fmt.Println("Processed: ", parsedEntry.Level, parsedEntry.Message)

	}

}

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
