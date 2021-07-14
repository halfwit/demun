package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/halfwit/demun/internal/command"
)

var (
	debug = flag.Bool("d", false, "Debug")
	port = flag.String("p", ":9997", "Port to listen on")
)

func incoming(li net.Listener, listen chan <- net.Conn) {
	for {
		conn, err := li.Accept()
		if err != nil {
			continue
		}
		listen <- conn
	}
}

func main() {
	flag.Parse()

	// Create a listener for any signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Start up our net listener
	li, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatal(err)
	}
	defer li.Close()

	listen := make(chan net.Conn, 50)
	go incoming(li, listen)

	cmd := command.NewCommand()
	if *debug {
		cmd.Logger = log.Printf
	}

	go cmd.Listen()

	for {
		select {
		case <-interrupt:
			return
		case conn := <-listen:
			go cmd.Handle(conn)
		}

	}
}
