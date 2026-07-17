package main

import (
	"fmt"
	"net"

	"github.com/kelvinjrosado/httpfromtcp/internal/request"
)

func main() {
	ls, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}
	defer ls.Close()

	for {

		conn, err := ls.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err)
			return
		}

		fmt.Println("Accepted a new connection")

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Error parsing incoming request:", err)
			return
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %v\n", req.RequestLine.Method)
		fmt.Printf("- Target: %v\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %v\n", req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf("- %v: %v\n", k, v)
		}
		fmt.Println("Body:")
		fmt.Println(string(req.Body))
		fmt.Println("Connection has been closed")
	}
}
