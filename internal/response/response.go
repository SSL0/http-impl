package response

import (
	"fmt"
	"io"
	"log/slog"
	"strconv"

	"github.com/SSL0/http-impl/internal/headers"
)

const (
	StatusOK                  = 200
	StatusBadRequset          = 400
	StatusInternalServerError = 500
)

const (
	httpVersion = "HTTP/1.1"
	CRLF        = "\r\n"
)

func StatusText(code int) string {
	switch code {
	case StatusOK:
		return "OK"
	case StatusBadRequset:
		return "Bad Request"
	case StatusInternalServerError:
		return "Internal Server Error"
	default:
		return ""
	}
}

func WriteStatusLine(w io.Writer, statusCode int) error {
	reasonPhrase := StatusText(statusCode)

	if reasonPhrase == "" {
		return fmt.Errorf("unknown status code: %d", statusCode)
	}

	// HTTP-version SP status-code SP [ reason-phrase ]
	statusLine := fmt.Sprintf("%s %d %s\r\n", httpVersion, statusCode, reasonPhrase)
	slog.Info("write status line", "status-line", statusLine)
	w.Write([]byte(statusLine))
	return nil
}

// Hardcoded
func GetDefaultHeaders(contentLength int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLength))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, h headers.Headers) error {
	data := []byte{}

	h.ForEach(func(k string, v string) {
		fieldLine := fmt.Sprintf("%s: %s\r\n", k, v)
		data = append(data, []byte(fieldLine)...)
	})
	data = append(data, []byte(CRLF)...)

	slog.Info("write status line", "data", data)
	n, err := w.Write(data)

	if err != nil {
		return err
	}

	if l := len(data); n != l {
		return fmt.Errorf("not all headers was written: written - %d, headers len - %d", n, l)
	}

	return nil
}
