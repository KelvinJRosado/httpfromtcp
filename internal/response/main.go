package response

import (
	"io"
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
		_, err := w.Write([]byte("HTTP/1.1 200 OK"))
		return err
	case Status400:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request"))
		return err
	case Status500:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error"))
		return err
	default:
		_, err := w.Write([]byte("HTTP/1.1 500 "))
		return err
	}
}
