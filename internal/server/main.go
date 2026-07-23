package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/kelvinjrosado/httpfromtcp/internal/request"
	"github.com/kelvinjrosado/httpfromtcp/internal/response"
)

// Contains the state of the server
type Server struct {
	closed   *atomic.Bool // if server is done accepting traffic
	listener net.Listener // network listerner to pull requests
	handler  Handler
}
type HandlerError struct {
	StatusCode int
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

// Starts listening for requests inside a goroutine.
func Serve(port int, handler Handler) (*Server, error) {
	lr, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return &Server{}, err
	}

	cl := atomic.Bool{}
	cl.Store(false)

	srv := Server{
		closed:   &cl,
		listener: lr,
		handler:  handler,
	}

	go func() {
		srv.listen()
	}()

	return &srv, nil
}

// Closes the listener and the server
func (s *Server) Close() error {
	s.closed.Store(true)

	err := s.listener.Close()
	if err != nil {
		return err
	}
	return nil
}

// Uses a loop to .Accept new connections as they come in, and handles each one in a new goroutine
func (s *Server) listen() {
	for {

		cn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}

			fmt.Printf("could not accept connection: %v\n", err)
			continue
		}

		go func() {
			s.handle(cn)
		}()

	}
}

// Handles a single connection
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		return
	}

	buf := bytes.Buffer{}

	he := s.handler(&buf, req)
	if he.Message != "" {
		_ = he.WriteHandlerError(conn)
		return
	}

	_ = response.WriteStatusLine(conn, response.Status200)

	hr := response.GetDefaultHeaders(0)
	_ = response.WriteHeaders(conn, hr)

	_ = response.WriteBody(conn, buf.Bytes())
}

func (he *HandlerError) WriteHandlerError(w io.Writer) error {
	_, err := fmt.Fprintf(w, "Failed with status %v: %v", he.StatusCode, he.Message)
	return err
}
