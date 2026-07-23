package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/kelvinjrosado/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	Status200 StatusCode = iota
	Status400
	Status500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case Status200:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return err
	case Status400:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return err
	case Status500:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		return err
	default:
		_, err := w.Write([]byte("HTTP/1.1 500 \r\n"))
		return err
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	hr := headers.NewHeaders()

	hr["content-length"] = strconv.Itoa(contentLen)
	hr["connection"] = "close"
	hr["content-type"] = "text/plain"

	return hr
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		line := fmt.Sprintf("%v: %v\r\n", k, v)

		_, err := w.Write([]byte(line))
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	return nil
}

func WriteBody(w io.Writer, body []byte) error {
	_, err := w.Write(body)
	return err
}
