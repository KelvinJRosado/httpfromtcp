package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	// Check for proxy path first
	path, isProxy := strings.CutPrefix(req.RequestLine.RequestTarget, "/httpbin")
	if isProxy {
		req.RequestLine.RequestTarget = path
		proxyHandler(w, req)
		return
	}

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

func proxyHandler(w *response.Writer, req *request.Request) {
	statusCode := response.Status200

	hs := response.GetDefaultHeaders(0)
	hs.UseChunked()
	hs["content-type"] = "application/json"

	if err := w.WriteStatusLine(statusCode); err != nil {
		return
	}

	if err := w.WriteHeaders(hs); err != nil {
		return
	}

	res, err := http.Get(fmt.Sprintf("https://httpbin.org%v", req.RequestLine.RequestTarget))
	if err != nil {
		return
	}
	defer res.Body.Close()

	buf := make([]byte, 32)
	for {
		readBytes, errRead := res.Body.Read(buf)

		if readBytes > 0 {
			_, errWrite := w.WriteChunkedBody(buf[:readBytes])
			if errWrite != nil {
				return
			}
		}

		if errRead != nil {
			if errRead == io.EOF {
				_, _ = w.WriteChunkedBodyDone()
			}

			return
		}

	}
}
