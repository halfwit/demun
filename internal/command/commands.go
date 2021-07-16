package command

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"regexp"
)

type Command struct {
	Logger	func(string, ...interface{})
	cmds	chan entry
}

var (
	add = []byte("add")
	list = []byte("list")
	remove = []byte("remove")
)

// This is fairly sloppy, and would do well with a refactoring
type entry struct {
	name	string
	data 	string
	tag	string
	regexp	*regexp.Regexp
	result	chan entry
}

func NewCommand() *Command {
	return &Command{
		Logger:	func(string, ...interface{}) {},
		cmds: make(chan entry),
	}
}

func (command *Command) Listen() {
	var data []entry

	command.Logger("Listening for commands")
	for cmd := range command.cmds {
		switch cmd.name {
		case "add":
			data = append(data, cmd)
		case "list":
			for _, item := range data {
				if item.tag == cmd.tag {
					cmd.result <- item
				}
			} 
			command.Logger("Sending list for %s\n", cmd.tag)
			close(cmd.result)
		case "remove":
			count := 0
			for _, item := range data {
				if ! cmd.regexp.MatchString(item.data) {
					data[count] = item
					count++
				}
			}
			data = data[:count]
		}
	}	
}

func (command *Command) Handle(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 1024*1024)
	if(!scanner.Scan()) {
		fmt.Fprintf(conn, "No input detected")
		return
	}

	// Scan in the first line
	target := scanner.Bytes()
	if bytes.HasPrefix(target, list) {
		result := make(chan entry)
		command.cmds <- entry{
			name: "list",
			tag: string(target[5:]),
			result: result,
		}
		writer := bufio.NewWriter(conn)
		for item := range result {
			writer.WriteString(item.data)
			writer.WriteByte('\n')
		}

		writer.Flush()
		return
	}

	if bytes.HasPrefix(target, remove) {
		rx, err := regexp.Compile(string(target[7:]))
		if err != nil {
			fmt.Fprintf(conn, "%s\n", "Invalid regex supplied")
			return
		}

		command.cmds <- entry{
			name: "remove",
			regexp: rx,
		}
	}

	if ! bytes.HasPrefix(target, add) {
		return
	}

	command.Logger("Adding entry for tag %s", target[4:])
	for scanner.Scan() {
		command.cmds <- entry{
			name: "add",
			tag: string(target[4:]),
			data: scanner.Text(),
		}
	}
}
