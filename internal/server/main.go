package server

import (
	"fmt"
	"net"
	"sync/atomic"
)

// Contains the state of the server
type Server struct {
	closed   *atomic.Bool  // if server is done accepting traffic
	listener *net.Listener // network listerner to pull requests
}

// Starts listening for requests inside a goroutine.
func Serve(port int) (*Server, error) {
	lr, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return &Server{}, nil
	}

	cl := atomic.Bool{}
	cl.Store(false)

	go func() {
		lr.Accept()
	}()

	return &Server{
		closed:   &cl,
		listener: &lr,
	}, nil
}

// Closes the listener and the server
func (s *Server) Close() error

// Uses a loop to .Accept new connections as they come in, and handles each one in a new goroutine
func (s *Server) listen()

// Handles a single connection
func (s *Server) handle(conn net.Conn)
