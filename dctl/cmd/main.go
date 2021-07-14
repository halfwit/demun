package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

//	"github.com/halfwit/demun/internal/command"
)

var (
	tag = flag.String("t", "path", "Set tag on data")
	srv = flag.String("s", "localhost", "Address of host")
	host = flag.String("r", "", "Host prefix for files")
	port = flag.String("p", "9997", "Port to connect via")
)

func list(conn net.Conn) {
	fmt.Fprintf(conn, "list %s\n", *tag)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
}

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
	case "list":
		list(conn)
	}

	conn.Close()
}
