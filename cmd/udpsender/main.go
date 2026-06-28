package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Failed to dial UDP:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading string", err)
			return
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Println("Error writing string:", err)
			return
		}
	}
}
