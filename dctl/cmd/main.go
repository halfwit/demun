package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"net"
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
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	defer conn.Close()

	// Add and list
	switch (flag.Arg(0)) {
	case "add":
		fmt.Fprintf(conn, "add %s\n", *tag)
		scanner := bufio.NewScanner(os.Stdin)
		buf := make([]byte, 1024*1024)
		scanner.Buffer(buf, 1024*1024)

		writer := bufio.NewWriter(conn)
		for scanner.Scan() {
			writer.WriteString(*host)
			writer.Write(scanner.Bytes() )
			writer.WriteByte('\n')
		}

		writer.Flush()
	case "list":
		fmt.Fprintf(conn, "list %s\n", *tag)
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	case "remove":
		if flag.NArg() <= 1 {
			log.Fatal("Remove must be supplied a regex to match entries")
		}
		fmt.Fprintf(conn, "remove %s\n", flag.Arg(1))
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}
}
