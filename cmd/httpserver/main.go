package main

import (
	"fmt"
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

	body := []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)

	hs := response.GetDefaultHeaders(len(body))
	hs["content-type"] = "text/html"

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		statusCode = response.Status400
		body = []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
	case "/myproblem":
		statusCode = response.Status500
		body = []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
	}

	hs["content-length"] = fmt.Sprintf("%d", len(body))

	if err := w.WriteStatusLine(statusCode); err != nil {
		return
	}

	if err := w.WriteHeaders(hs); err != nil {
		return
	}

	_, _ = w.WriteBody(body)
}
