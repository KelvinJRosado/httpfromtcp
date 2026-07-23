package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelvinjrosado/httpfromtcp/internal/request"
	"github.com/kelvinjrosado/httpfromtcp/internal/response"
	"github.com/kelvinjrosado/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, myHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func myHandler(w io.Writer, req *request.Request) *server.HandlerError {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{
			StatusCode: response.Status400,
			Message:    "Your problem is not my problem",
		}
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{
			StatusCode: response.Status500,
			Message:    "Woopsie, my bad",
		}
	}

	_, err := w.Write([]byte("All good, frfr"))
	if err != nil {
		return &server.HandlerError{
			StatusCode: response.Status500,
			Message:    fmt.Sprintf("error writing from handler: %v", err),
		}
	}

	return nil
}
