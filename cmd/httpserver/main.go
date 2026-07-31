package main

import (
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

func myHandler(w *response.Writer, req *request.Request) {
	statusCode := response.Status200
	body := []byte("All good, frfr")

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		statusCode = response.Status400
		body = []byte("Your problem is not my problem")
	case "/myproblem":
		statusCode = response.Status500
		body = []byte("Woopsie, my bad")
	}

	if err := w.WriteStatusLine(statusCode); err != nil {
		return
	}

	h := response.GetDefaultHeaders(len(body))
	if err := w.WriteHeaders(h); err != nil {
		return
	}

	_, _ = w.WriteBody(body)
}
