package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
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

	s := string(b)

	lines := strings.Split(s, "\r\n")

	rl, err := parseRequestLine(lines[0])
	if err != nil {
		return &req, err
	}

	req.RequestLine = *rl

	return &req, nil
}

func parseRequestLine(line string) (*RequestLine, error) {
	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		return nil, fmt.Errorf("request line muct have exactly 3 parts. Received: %v", len(parts))
	}

	method := parts[0]
	path := parts[1]
	version := strings.ReplaceAll(parts[2], "HTTP/", "")

	// Validate method
	if method != strings.ToUpper(method) {
		return nil, fmt.Errorf("method must be all caps. Received: %v", method)
	}

	// Validate path

	// Validate version
	if version != "1.1" {
		return nil, fmt.Errorf("only HTTP/1.1 is supported. Received: %v", version)
	}

	rl := RequestLine{
		HttpVersion:   version,
		RequestTarget: path,
		Method:        method,
	}

	return &rl, nil
}
