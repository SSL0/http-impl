package server

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync/atomic"

	"github.com/SSL0/http-impl/internal/request"
	"github.com/SSL0/http-impl/internal/response"
)

type HandlerError struct {
	StatusCode int
	Message    string
}

func (e *HandlerError) Write(w io.Writer) error {
	data := fmt.Sprintf("%d %s", e.StatusCode, e.Message)
	n, err := w.Write([]byte(data))

	if n != len(data) {
		return fmt.Errorf("failed to write full handler error data: %s", data)
	}

	return err
}

type HandlerFunc func(w *response.Writer, req *request.Request)

type Server struct {
	listener *net.Listener
	handler  HandlerFunc
	closed   atomic.Bool
}

func newServer(l *net.Listener, f HandlerFunc) *Server {
	return &Server{
		listener: l,
		handler:  f,
		closed:   atomic.Bool{},
	}
}

func ListenAndServe(port uint16, f HandlerFunc) (*Server, error) {
	address := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", address)

	if err != nil {
		return nil, err
	}

	server := newServer(&l, f)

	go server.Serve()

	return server, nil
}

func (s *Server) Serve() {
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
	defer conn.Close()

	req, err := request.RequestFromReader(conn)

	if err != nil {
		slog.Error("failed to get request from client", "context_error", err)

		errMsg := fmt.Sprintf(
			"%d %s",
			response.StatusBadRequset,
			response.StatusText(response.StatusBadRequset),
		)

		conn.Write([]byte(errMsg))
		return
	}

	rWriter := response.NewWriter(conn)
	s.handler(rWriter, req)
}

func (s *Server) Close() {
	s.closed.Store(true)
}
