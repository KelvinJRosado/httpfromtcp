package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// convert to string for utility
	rawLine := string(data)

	// If no CRLF present, we need more data
	if !strings.Contains(rawLine, crlf) {
		return 0, false, nil
	}

	// If CRLF is at start, we are finished
	if strings.HasPrefix(rawLine, crlf) {
		return 0, true, nil
	}

	// Only process 1st header (1 at a time)
	line := strings.Split(rawLine, crlf)[0]

	// Extract the name and value
	fieldLine := strings.SplitN(line, ":", 2)

	if len(fieldLine) != 2 {
		return 0, false, fmt.Errorf("missing field name and/or value. Received: %v", fieldLine)
	}

	fieldName := strings.ToLower(fieldLine[0])

	// Ensure name is valid
	if strings.HasPrefix(fieldName, " ") || strings.HasSuffix(fieldName, " ") {
		return 0, false, fmt.Errorf("header name cannot have spaces at start or end. Received: %v", fieldName)
	}

	// Remove optional whitespace
	fieldValue := strings.TrimSpace(fieldLine[1])

	// Update map
	h[fieldName] = fieldValue

	return (len(line) + len(crlf)), false, nil
}
