package response

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/kelvinjrosado/httpfromtcp/internal/headers"
)

type (
	StatusCode   int
	writerStatus int
)

type Writer struct {
	writer      io.Writer
	writerState writerStatus
}

const (
	writerStatusInit writerStatus = iota
	writerStatusWroteStatusLine
	writerStatusWroteHeaders
	writerStatusErrored
)

const (
	Status200 StatusCode = iota
	Status400
	Status500
)

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer, writerState: writerStatusInit}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	hr := headers.NewHeaders()

	hr["content-length"] = strconv.Itoa(contentLen)
	hr["connection"] = "close"
	hr["content-type"] = "text/plain"

	return hr
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writerStatusInit {
		return errors.New("status line must be written first")
	}

	var err error

	switch statusCode {
	case Status200:
		_, err = w.writer.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case Status400:
		_, err = w.writer.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case Status500:
		_, err = w.writer.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		_, err = w.writer.Write([]byte("HTTP/1.1 500 \r\n"))
	}

	if err != nil {
		w.writerState = writerStatusErrored
	} else {
		w.writerState = writerStatusWroteStatusLine
	}

	return err
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.writerState != writerStatusWroteStatusLine {
		return errors.New("headers must be written second")
	}

	for k, v := range h {
		line := fmt.Sprintf("%v: %v\r\n", k, v)

		_, err := w.writer.Write([]byte(line))
		if err != nil {
			w.writerState = writerStatusErrored
			return err
		}
	}

	_, err := w.writer.Write([]byte("\r\n"))
	if err != nil {
		w.writerState = writerStatusErrored
		return err
	}

	w.writerState = writerStatusWroteHeaders
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writerStatusWroteHeaders {
		return 0, errors.New("body must be written last")
	}

	written, err := w.writer.Write(p)
	if err != nil {
		w.writerState = writerStatusErrored
	}

	return written, err
}
