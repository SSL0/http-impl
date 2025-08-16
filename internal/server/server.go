package server

import (
	"fmt"
	"log/slog"
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
	for !s.closed.Load() {
		conn, err := (*s.listener).Accept()

		if err != nil {
			slog.Error("failed to listen", "context_error", err)
		}

		slog.Info("conn successfully accepted", "remote_client_ip", conn.RemoteAddr().String())

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	response := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!\n")
	_, err := conn.Write(response)
	if err != nil {
		slog.Error("failed to write to conn", "context_error", err)
	}

	err = conn.Close()

	if err != nil {
		slog.Error("failed to close conn", "context_error", err)
	}
}

func (s *Server) Close() {
	s.closed.Store(true)
}
