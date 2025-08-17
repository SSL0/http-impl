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
	StatusNotFound            = 404
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
	case StatusNotFound:
		return "Not Found"
	case StatusInternalServerError:
		return "Internal Server Error"
	default:
		return ""
	}
}

type writerState int

const (
	WritingStatusLine writerState = iota
	WritingHeaders    writerState = iota
	WritingBody       writerState = iota
)

type Writer struct {
	writer io.Writer
	state  writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w, state: WritingStatusLine}
}

func (w *Writer) WriteStatusLine(statusCode int) error {
	if w.state != WritingStatusLine {
		return fmt.Errorf("failed to write status line, writer state is different")
	}
	err := writeStatusLine(w.writer, statusCode)
	if err != nil {
		return err
	}
	w.state = WritingHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != WritingHeaders {
		return fmt.Errorf("failed to write status line, writer state is different")
	}

	err := writeHeaders(w.writer, headers)
	if err != nil {
		return err
	}
	w.state = WritingBody

	return nil
}

func (w *Writer) WriteBody(p []byte) error {
	if w.state != WritingBody {
		return fmt.Errorf("failed to write status line, writer state is different")
	}

	err := writeBody(w.writer, p)
	if err != nil {
		return err
	}
	w.state = WritingBody

	return nil
}

func writeStatusLine(w io.Writer, statusCode int) error {
	reasonPhrase := StatusText(statusCode)

	if reasonPhrase == "" {
		return fmt.Errorf("unknown status code: %d", statusCode)
	}

	// HTTP-version SP status-code SP [ reason-phrase ]
	statusLine := fmt.Sprintf("%s %d %s\r\n", httpVersion, statusCode, reasonPhrase)
	slog.Info("wrote status line", "status-line", statusLine)
	w.Write([]byte(statusLine))
	return nil
}

func writeHeaders(w io.Writer, h headers.Headers) error {
	data := []byte{}

	h.ForEach(func(k string, v string) {
		fieldLine := fmt.Sprintf("%s: %s\r\n", k, v)
		data = append(data, []byte(fieldLine)...)
	})
	data = append(data, []byte(CRLF)...)

	slog.Info("wrote headers", "data", data)
	n, err := w.Write(data)

	if err != nil {
		return err
	}

	if l := len(data); n != l {
		return fmt.Errorf("failed to write all headers data: written - %d, headers len - %d", n, l)
	}

	return nil
}

func writeBody(w io.Writer, body []byte) error {
	n, err := w.Write(body)

	if err != nil {
		return nil
	}
	slog.Info("wrote body", "data", body)

	if l := len(body); n != l {
		return fmt.Errorf("failed to write all body: written - %d, headers len - %d", n, l)
	}
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
