package server

import (
	"fmt"
	"net"
	"sync/atomic"
)

// Contains the state of the server
type Server struct {
	closed   *atomic.Bool // if server is done accepting traffic
	listener net.Listener // network listerner to pull requests
}

// Starts listening for requests inside a goroutine.
func Serve(port int) (*Server, error) {
	lr, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return &Server{}, err
	}

	cl := atomic.Bool{}
	cl.Store(false)

	srv := Server{
		closed:   &cl,
		listener: lr,
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

	data := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\n"

	_, err := conn.Write([]byte(data))
	if err != nil {
		fmt.Printf("could not write response: %v\n", err)
	}
}
