package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SSL0/http-impl/internal/request"
	"github.com/SSL0/http-impl/internal/response"
	"github.com/SSL0/http-impl/internal/server"
)

const port = 42069

const htmlBadRequest = `<html>
  <head>
	<title>400 Bad Request</title>
  </head>
  <body>
	<h1>Bad Request</h1>
	<p>Your request honestly kinda sucked.</p>
  </body>
</html>
`

const htmlInternalServerError = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`

const htmlOK = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`

const htmlNotFound = `<html>
  <head>
    <title>404 Not Found</title>
  </head>
  <body>
    <h1>Not Found</h1>
    <p>Unknown resource</p>
  </body>
</html>
`

func serverHandle(resWriter *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		resWriter.WriteStatusLine(response.StatusBadRequset)
		h := response.GetDefaultHeaders(len(htmlBadRequest))
		h.Change("Content-Type", "text/html")
		resWriter.WriteHeaders(h)
		resWriter.WriteBody([]byte(htmlBadRequest))
	case "/myproblem":
		resWriter.WriteStatusLine(response.StatusInternalServerError)
		h := response.GetDefaultHeaders(len(htmlInternalServerError))
		h.Change("Content-Type", "text/html")
		resWriter.WriteHeaders(h)
		resWriter.WriteBody([]byte(htmlInternalServerError))
	case "/correct":
		resWriter.WriteStatusLine(response.StatusOK)
		h := response.GetDefaultHeaders(len(htmlOK))
		h.Change("Content-Type", "text/html")
		resWriter.WriteHeaders(h)
		resWriter.WriteBody([]byte(htmlOK))
	default:
		resWriter.WriteStatusLine(response.StatusNotFound)
		h := response.GetDefaultHeaders(len(htmlNotFound))
		h.Change("Content-Type", "text/html")
		resWriter.WriteHeaders(h)
		resWriter.WriteBody([]byte(htmlNotFound))
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
