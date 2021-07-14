package command

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
)

type Command struct {
	Type string
	Tag string
	Entry *Entry
	Result chan *Entry 
}

type Entry struct {
	Title	string
	Tag string
	Payload	string
	Option	string
}

func Listen(commands chan Command) {
	var data []*Entry
	for cmd := range commands {
		switch cmd.Type {
		case "add":
			data = append(data, cmd.Entry) 
		case "list":
			for _, item := range data {
				if item.Tag == cmd.Tag {
					cmd.Result <- item
				}
			} 
		}
	}	
}

func Handle(commands chan Command, conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	if(!scanner.Scan()) {
		fmt.Fprintf(conn, "No input detected")
		return
	}
	// Scan in the first line
	target := scanner.Bytes()
	if bytes.Contains(target, []byte("list")) {
		result := make(chan *Entry)
		commands <- Command{
			Type: "list",
			Tag: string(target[4:]),
			Result: result,
		}
		for entry := range result {
			fmt.Fprintf(conn, "%s - %s - %s\n", entry.Title, entry.Payload, entry.Option)
		}
		return
	}	

	if ! bytes.Contains(target, []byte("add")) {
		return
	}

	for scanner.Scan() {
		ln := scanner.Text()
	 	entry := &Entry{
			Title: "test",
			Tag: string(target[3:]),
			Payload: ln,
		}
		commands <- Command{
			Type: "add",
			Entry: entry,
		}
	}
}
