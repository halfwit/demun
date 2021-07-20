package command

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"regexp"
	"sort"
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

type entries []entry

func (e entries) Len()	int		{ return len(e) }
func (e entries) Swap(i, j int)		{ e[i], e[j] = e[j], e[i] }
func (e entries) Less(i, j int) bool	{ return e[i].data < e[j].data } 
func NewCommand() *Command {
	return &Command{
		Logger:	func(string, ...interface{}) {},
		cmds: make(chan entry),
	}
}
func (command *Command) Listen() {
	var data entries 

	command.Logger("Listening for commands")
	for cmd := range command.cmds {
		switch cmd.name {
		case "add":
			data = append(data, cmd)
			sort.Sort(data)
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
		writer := bufio.NewWriterSize(conn, 1024*1024)
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
