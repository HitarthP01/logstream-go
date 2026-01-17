package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	// Start listening on TCP port 9000
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server started on port 9000....")

	// Accept connections in a loop, which is infinite loop here
	for {
		conn, err := listener.Accept() //Accept() blocks (waits) until a client connects
		/*
			┌────────────┐                      ┌────────────┐
			│   Client   │                      │   Server   │
			│            │                      │            │
			│ net.Dial() │ ──── connects ────▶  │  Accept()  │
			│     │      │                      │     │      │
			│     ▼      │                      │     ▼      │
			│   conn     │ ◄═══════════════════▶│   conn     │
			└────────────┘    TCP Connection    └────────────┘

		*/
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
		// concurrency with goroutines, with the keyterm 'go'
		//explain go handleConnection(conn){}
		// The 'go' keyword is used to start a new goroutine, which is a lightweight thread managed by the Go runtime.
		// By prefixing the function call with 'go', the program can handle multiple connections concurrently without blocking the main thread.
		// This allows the server to efficiently manage multiple clients at the same time.

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {

		line := scanner.Text()
		fmt.Println("Received:", line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading:", err)
	}
	fmt.Println("Client disconnected: ", conn.RemoteAddr())
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
