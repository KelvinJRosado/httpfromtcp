package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Status      int // 0=init, 1=done
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)

	req := Request{
		Status: 0,
	}

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
		return nil, 0, nil
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

func (r *Request) parse(data []byte) (int, error) {
	var rl *RequestLine
	var read int
	var err error

	switch r.Status {
	case 0:
		rl, read, err = parseRequestLine(data)
	case 1:
		return read, fmt.Errorf("trying to read done but request state is done")
	default:
		return read, fmt.Errorf("unknown request state")
	}

	if err != nil {
		return read, err
	}

	if read == 0 {
		return read, nil
	}

	r.RequestLine = *rl
	r.Status = 1

	return read, nil
}
