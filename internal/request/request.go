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

const bufferSize = 8 // Number of bytes to read at a time

func RequestFromReader(reader io.Reader) (*Request, error) {
	readToIndex := 0 // How much we have read so far
	req := Request{
		Status: 0,
	}

	buf := make([]byte, bufferSize)

	for req.Status != 1 {
		// Grow buffer if needed
		if readToIndex == cap(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		// Read at least
		read, err := io.ReadAtLeast(reader, buf[readToIndex:], 1)
		if err == io.EOF {
			req.Status = 1
			break
		}

		readToIndex += read

		if err != nil {
			return &req, err
		}

		pd, err := req.parse(buf[:readToIndex])
		if err != nil {
			return &req, err
		}

		if pd == 0 {
			continue
		}

		newBuf := make([]byte, len(buf))
		copy(newBuf, buf[pd:readToIndex])
		buf = newBuf

		readToIndex -= pd

	}

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
		return nil, len(line), fmt.Errorf("request line must have exactly 3 parts. Received: %v, %v", len(parts), parts)
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

	return &rl, len(line + "\r\n"), nil
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
