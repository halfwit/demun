package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var (
	tag = flag.String("t", "path", "Set tag on data")
	srv = flag.String("s", "localhost", "Address of host")
	host = flag.String("r", "", "Host prefix for files")
	port = flag.String("p", "9997", "Port to connect via")
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%s", *srv, *port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	// Make sure we don't block for long
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Add and list
	switch (flag.Arg(0)) {
	case "add":
		fmt.Fprintf(conn, "add %s\n", *tag)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Fprintf(conn, "%s\n", scanner.Text())
		}
	case "list":
		fmt.Fprintf(conn, "list %s\n", *tag)
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}

	conn.Close()
}
