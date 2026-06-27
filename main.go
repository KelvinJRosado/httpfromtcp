package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	// Open file for reading
	f, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	ch := getLinesChannel(f)

	for st := range ch {
		fmt.Printf("read: %s\n", st)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	messages := make(chan string)

	go func() {
		// Vars for keeping track of progress
		curr := ""
		for {

			// Buffer to hold the next 8 bytes we read from file
			buf := make([]byte, 8)

			// Read the next 8 bytes and save to buffer
			_, err := io.ReadAtLeast(f, buf, 1)
			isEOF := errors.Is(err, io.EOF) // Detect end of file
			if err != nil && !isEOF {
				fmt.Println("Unexpected error reading file:", err)
				break
			}

			//tmp := ""

			// Populate our string holding a full line
			for _, r := range string(buf) {
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

		close(messages)
	}()
	return messages
}
