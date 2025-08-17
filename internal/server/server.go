package server

import (
	"fmt"
	"log/slog"
	"net"
	"sync/atomic"

	"github.com/SSL0/http-impl/internal/response"
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
	err := response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		slog.Error("failed to write status line", "context_error", err)
	}

	err = response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	if err != nil {
		slog.Error("failed to write headers", "context_error", err)
	}

	err = conn.Close()
	if err != nil {
		slog.Error("failed to close conn", "context_error", err)
	}
}

func (s *Server) Close() {
	s.closed.Store(true)
}
