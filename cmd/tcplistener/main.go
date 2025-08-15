package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func(ch chan string) {
		currentLine := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}

			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				currentLine += string(data[:i])
				data = data[i+1:]
				ch <- currentLine
				currentLine = ""
			}
			currentLine += string(data)
		}
		if len(currentLine) != 0 {
			ch <- currentLine
		}
		close(ch)
	}(ch)

	return ch
}

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

		ch := getLinesChannel(conn)
		for {
			v, ok := <-ch
			if !ok {
				break
			}
			fmt.Printf("%s\n", v)
		}
		err = conn.Close()
		if err != nil {
			log.Fatalf("failed to close connection %v", err)
		}
		log.Printf("conn successfully closed\n")

	}
}
