package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener *net.Listener
	closed   atomic.Bool
}

func newServer(listener *net.Listener) *Server {
	return &Server{listener: listener, closed: atomic.Bool{}}
}

func Serve(port uint16) (*Server, error) {
	address := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", address)

	if err != nil {
		return nil, err
	}

	server := newServer(&l)

	go server.listen()

	return server, nil
}

func (s *Server) listen() {
	s.closed.Store(false)
	for {
		if s.closed.Load() {
			return
		}

		conn, err := (*s.listener).Accept()

		if err != nil {
			log.Fatalf("failed to accept connection %v", err)
		}

		log.Printf("conn successfully accepted\n")

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	response := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!\n")
	conn.Write(response)
	conn.Close()
}

func (s *Server) Close() {
	s.closed.Store(true)
}
