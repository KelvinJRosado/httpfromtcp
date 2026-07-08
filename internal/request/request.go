package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Status      int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)

	req := Request{}

	if err != nil {
		return &req, err
	}

	rl, _, err := parseRequestLine(b)
	if err != nil {
		return &req, err
	}

	req.RequestLine = *rl

	return &req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	s := string(data)

	lines := strings.Split(s, "\r\n")

	if len(lines) < 2 {
		return nil, 0, fmt.Errorf("no CRLF found in data")
	}

	line := lines[0]

	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		return nil, len(line), fmt.Errorf("request line muct have exactly 3 parts. Received: %v", len(parts))
	}

	method := parts[0]
	path := parts[1]
	version := strings.ReplaceAll(parts[2], "HTTP/", "")

	// Validate method
	if method != strings.ToUpper(method) {
		return nil, len(line), fmt.Errorf("method must be all caps. Received: %v", method)
	}

	// Validate path

	// Validate version
	if version != "1.1" {
		return nil, len(line), fmt.Errorf("only HTTP/1.1 is supported. Received: %v", version)
	}

	rl := RequestLine{
		HttpVersion:   version,
		RequestTarget: path,
		Method:        method,
	}

	return &rl, len(line), nil
}
