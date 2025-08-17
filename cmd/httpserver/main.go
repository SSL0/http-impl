package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SSL0/http-impl/internal/request"
	"github.com/SSL0/http-impl/internal/server"
)

const port = 42069

func serverHandle(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			StatusCode: 400,
			Message:    "Your problem is not my problem\n",
		}
	case "/myproblem":
		return &server.HandlerError{
			StatusCode: 500,
			Message:    "My bad\n",
		}
	default:
		w.Write([]byte("All good\n"))
		return nil
	}
}

func main() {
	server, err := server.ListenAndServe(port, serverHandle)

	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	defer server.Close()

	log.Printf("server started on port %d\n", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Printf("found signal to stop server")
}
