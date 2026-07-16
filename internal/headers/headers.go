package headers

import (
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

var regexHeaderName = regexp.MustCompile("^[A-Za-z0-9!#$%&'*+\\-.^_`|~]+$")

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
		return len(crlf), true, nil
	}

	// Only process 1st header (1 at a time)
	line, _, _ := strings.Cut(rawLine, crlf)

	// Extract the name and value
	div := strings.Index(line, ":")

	if div < 1 {
		return 0, false, fmt.Errorf("no : inside field-line. Received: %v", line)
	}

	fieldName := strings.ToLower(line[:div])

	// Ensure name is valid
	if fieldName != strings.TrimSpace(fieldName) {
		return 0, false, fmt.Errorf("header name cannot have whitespace at start or end. Received: %v", fieldName)
	}

	if !regexHeaderName.MatchString(fieldName) {
		return 0, false, fmt.Errorf("header name contains invalid characters. Received: %v", fieldName)
	}

	// Remove optional whitespace
	fieldValue := strings.TrimSpace(line[div+1:])

	// Update map

	val, found := h[fieldName]
	if found {
		h[fieldName] = val + ", " + fieldValue
	} else {
		h[fieldName] = fieldValue
	}

	return (len(line) + len(crlf)), false, nil
}

func (h Headers) Get(key string) (string, bool) {
	lower := strings.ToLower(key)

	val, ok := h[lower]
	return val, ok
}
