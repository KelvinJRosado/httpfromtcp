package request

import (
	"fmt"
	"io"
	"strconv"
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
	statusRequestStateParsingHeaders = 1
	statusDone                       = 2
	statusParsingBody                = 3
)

// Pulls data from an IO reader until we can parse a request
func RequestFromReader(reader io.Reader) (*Request, error) {
	readToIndex := 0 // How much we have read so far
	req := Request{
		Status: statusInit,
	}

	// Initialize our request data
	req.Headers = headers.NewHeaders()
	req.Body = []byte{}

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

		// We shouldn't hit EOF. Our status should end naturally before we hit EOF
		if err == io.EOF {
			return &req, io.ErrUnexpectedEOF
		}

		// Keep track of how much data we've read
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

	// Keep parsing until done processing the data
	for r.Status != statusDone {
		// Only parse is request is in init status. Error otherwise
		switch r.Status {
		// Start with request line parsing
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

			// Parse headers after request line
			r.Status = statusRequestStateParsingHeaders
		case statusRequestStateParsingHeaders:
			// Parse each header
			for {
				read, done, err := r.Headers.Parse(data[totalBytesParsed:])
				if err != nil {
					return totalBytesParsed, err
				}

				totalBytesParsed += read

				// Move on to body parsing after headers
				if done {
					r.Status = statusParsingBody
					break
				}

				// If not enough data, wait until next delivery
				if read == 0 {
					return totalBytesParsed, nil
				}

			}

		case statusParsingBody:
			// Make sure we only pull data based on reported length
			lenStr, found := r.Headers.Get("content-length")
			if !found {
				r.Status = statusDone
				break
			}
			length, err := strconv.Atoi(lenStr)
			if err != nil {
				return totalBytesParsed, err
			}

			// Check how many bytes left in the given data are for the body
			bytesLeft := len(data) - totalBytesParsed

			// Add the next body bytes to our current body
			r.Body = append(r.Body, data[totalBytesParsed:]...)
			totalBytesParsed += bytesLeft

			// Error if too much data
			if len(r.Body) > length {
				return totalBytesParsed, fmt.Errorf("data passed was longer than stated content-length. Expected: %v. Received: %v", length, len(r.Body))
			}

			// Done if data matches
			if len(r.Body) == length {
				r.Status = statusDone
				return totalBytesParsed, nil
			}

			// Of no more data to process, wait until next delivery
			if bytesLeft == 0 {
				return totalBytesParsed, nil
			}
		case statusDone:
			return totalBytesParsed, fmt.Errorf("trying to read done but request state is done")
		default:
			return totalBytesParsed, fmt.Errorf("unknown request state")
		}
	}
	return totalBytesParsed, nil
}
