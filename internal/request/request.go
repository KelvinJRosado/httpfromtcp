package request

import (
	"fmt"
	"io"
	"strings"

	"github.com/kelvinjrosado/httpfromtcp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Status      int
	Headers     headers.Headers
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	bufferSize = 8 // Number of bytes to read at a time
	crlf       = "\r\n"
	// Potential states our parsing can be in
	statusInit                       = 0
	statusDone                       = 2
	statusRequestStateParsingHeaders = 1
)

// Pulls data from an IO reader until we can parse a request
func RequestFromReader(reader io.Reader) (*Request, error) {
	readToIndex := 0 // How much we have read so far
	req := Request{
		Status: statusInit,
	}

	req.Headers = headers.NewHeaders()

	// Buffer to store bytes being processed
	buf := make([]byte, bufferSize)

	// Keep reading bytes until we have enough
	for req.Status != statusDone {
		// Grow buffer if needed
		if readToIndex == cap(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		// Read at least 1 byte and add to buffer
		numRead, err := io.ReadAtLeast(reader, buf[readToIndex:], 1)
		if err == io.EOF {
			req.Status = statusDone
			break
		}
		readToIndex += numRead
		if err != nil {
			return &req, err
		}

		// Attempt to parse buffered data
		numParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return &req, err
		}

		// If nothing to parse, keep looping
		if numParsed == 0 {
			continue
		}

		// Clean up buffer, removing parsed bytes
		copy(buf, buf[numParsed:readToIndex])
		readToIndex -= numParsed

	}

	return &req, nil
}

// Parse a given buffer of data, and extract an HTTP request line
func parseRequestLine(data []byte) (*RequestLine, int, error) {
	// Convert to string for convenience functions
	s := string(data)

	// Split on CRLF to split by HTTP spec standards
	lines := strings.Split(s, crlf)

	// Check if not enough data to parse
	if len(lines) < 2 {
		return nil, 0, nil
	}

	// Request line is 1st line in the request
	line := lines[0]

	// Split into the 3 parts required by the protocol
	parts := strings.Split(line, " ")

	// If we don't match spec, error
	if len(parts) != 3 {
		return nil, len(line), fmt.Errorf("request line must have exactly 3 parts. Received: %v, %v", len(parts), parts)
	}

	// Extract each part of the request line
	method := parts[0]
	path := parts[1]
	version := strings.ReplaceAll(parts[2], "HTTP/", "")

	// Validate method
	if method != strings.ToUpper(method) {
		return nil, len(line), fmt.Errorf("method must be all caps. Received: %v", method)
	}

	// No path validation for now

	// Validate version is 1.1, as per assignment
	if version != "1.1" {
		return nil, len(line), fmt.Errorf("only HTTP/1.1 is supported. Received: %v", version)
	}

	// Creat struct with parsed data
	rl := RequestLine{
		HttpVersion:   version,
		RequestTarget: path,
		Method:        method,
	}

	return &rl, len(line + crlf), nil
}

// Parse the provided buffer to populate the request
func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0

	for r.Status != statusDone {
		// Only parse is request is in init status. Error otherwise
		switch r.Status {
		case statusInit:
			rl, read, err := parseRequestLine(data[totalBytesParsed:])
			if err != nil {
				return read, err
			}

			if read == 0 {
				return read, nil
			}

			r.RequestLine = *rl
			totalBytesParsed += read

			r.Status = statusRequestStateParsingHeaders
		case statusRequestStateParsingHeaders:

			for {
				read, done, err := r.Headers.Parse(data[totalBytesParsed:])
				if err != nil {
					return totalBytesParsed, err
				}

				totalBytesParsed += read

				if done {
					r.Status = statusDone
					break
				}

				if read == 0 {
					return totalBytesParsed, nil
				}

			}

		case statusDone:
			return totalBytesParsed, fmt.Errorf("trying to read done but request state is done")
		default:
			return totalBytesParsed, fmt.Errorf("unknown request state")
		}
	}
	return totalBytesParsed, nil
}
