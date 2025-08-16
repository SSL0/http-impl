package main

import (
	"fmt"
	"log"
	"net"

	"github.com/SSL0/http-impl/internal/request"
)

func main() {
	l, err := net.Listen("tcp", "localhost:42069")

	if err != nil {
		log.Fatalf("failed to listen address %v", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("failed to accept connection %v", err)
		}

		log.Printf("conn successfully accepted\n")

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("failed make request from reader %v", err)
		}
		err = conn.Close()
		if err != nil {
			log.Fatalf("failed to close connection %v", err)
		}

		fmt.Printf("Request line:\n")
		fmt.Printf(
			"- Method: %s\n- Target: %s\n- Version: %s\n",
			req.RequestLine.Method,
			req.RequestLine.RequestTarget,
			req.RequestLine.HttpVersion,
		)

		fmt.Printf("Headers:\n")
		for k, v := range req.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
	}
}
