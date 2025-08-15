package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("failed to resolve udp addr %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("failed to dial udp%v", err)
	}
	defer conn.Close()

	r := bufio.NewReader(os.Stdin)

	for {
		_, err := conn.Write([]byte{'>'})
		if err != nil {
			log.Fatalf("failed write to udp conn %v", err)
		}

		str, err := r.ReadString('\n')
		if err != nil {
			log.Fatalf("failed write to read string from stdin %v", err)
		}

		_, err = conn.Write([]byte(str))
		if err != nil {
			log.Fatalf("failed write to write to udp conn %v", err)
		}
	}
}
