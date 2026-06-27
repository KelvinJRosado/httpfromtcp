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

	// Vars for keeping track of progress in file
	off := int64(0)
	curr := ""
	for {

		// Buffer to hold the next 8 bytes we read from file
		buf := make([]byte, 8)

		// Read the next 8 bytes and save to buffer
		n, err := f.ReadAt(buf, off)
		isEOF := errors.Is(err, io.EOF) // Detect end of file
		if err != nil && !isEOF {
			fmt.Println("Unexpected error reading file:", err)
			return
		}

		//tmp := ""

		// Populate our string holding a full line
		for _, r := range []rune(string(buf)) {
			if r == '\n' {
				fmt.Printf("read: %s\n", curr)
				curr = ""
			} else {
				curr += string(r)
			}
		}

		// Break from loop once we hit end of file
		if isEOF {
			break
		}

		// Increment byte pointer in file
		off += int64(n)
	}
}
