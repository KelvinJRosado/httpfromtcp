package main

import (
	"errors"
	"fmt"
	"io"
	"net"
)

func main() {
	ls, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}
	defer ls.Close()

	for {

		conn, err := ls.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err)
			return
		}

		fmt.Println("Accepted a new connection")

		ch := getLinesChannel(conn)

		for st := range ch {
			fmt.Printf("%s\n", st)
		}

		fmt.Println("Connection has been closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	messages := make(chan string)

	go func() {
		defer f.Close()
		// Vars for keeping track of progress
		curr := ""
		for {

			// Buffer to hold the next 8 bytes we read from file
			buf := make([]byte, 8)

			// Read the next 8 bytes and save to buffer
			numBytes, err := io.ReadAtLeast(f, buf, 1)
			isEOF := errors.Is(err, io.EOF) // Detect end of file
			if err != nil && !isEOF {
				fmt.Println("Unexpected error reading file:", err)
				break
			}

			// Populate our string holding a full line
			for _, r := range string(buf[:numBytes]) {
				if r == '\n' {
					messages <- curr
					curr = ""
				} else {
					curr += string(r)
				}
			}

			// Break from loop once we hit end of file
			if isEOF {
				break
			}

		}

		// Flush remaining data
		if curr != "" {
			messages <- curr
		}

		close(messages)
	}()
	return messages
}
