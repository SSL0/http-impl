package request

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
	"unicode"
)

const (
	bufferSize = 8
	crlf       = "\r\n"
)

var ErrMalformedRequestLine = errors.New("malformed request-line")
var ErrMissingSeparator = errors.New("missing separator")
var ErrInvalidMethodName = errors.New("incorrect method name")
var ErrIncorrectRequestTarget = errors.New("incorrect request target")
var ErrInvalidProtocolOrVersion = errors.New("incorrect protocol name or protocol version, only HTTP/1.1 allowed")

type parserState string

const (
	StateInitialized parserState = "initialized"
	StateDone        parserState = "done"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	RequestLine RequestLine
	state       parserState
}

func NewRequest() *Request {
	return &Request{state: StateInitialized}
}
func (r *Request) done() bool {
	return r.state == StateDone
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.state {
		case StateInitialized:
			rl, n, err := parseRequestLine(data[read:])

			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n

			r.state = StateDone
		case StateDone:
			break outer
		default:
			return 0, errors.New("unknown parser state")
		}
	}
	return read, nil
}

func parseRequestLine(s []byte) (*RequestLine, int, error) {
	if s[0] == ' ' {
		return nil, 0, ErrMalformedRequestLine
	}

	idx := bytes.Index(s, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := s[:idx]

	startLineParts := bytes.Split(startLine, []byte{' '})
	if len(startLineParts) != 3 {
		return nil, 0, ErrMalformedRequestLine
	}

	method := string(startLineParts[0])
	reqTarget := string(startLineParts[1])
	protocol := string(startLineParts[2])

	if !isUpperAndLetters(method) {
		return nil, 0, ErrInvalidMethodName
	}

	if _, err := url.Parse(reqTarget); err != nil {
		return nil, 0, ErrIncorrectRequestTarget
	}

	protocolParts := strings.Split(protocol, "/")

	if len(protocolParts) != 2 || protocolParts[0] != "HTTP" || protocolParts[1] != "1.1" {
		return nil, 0, ErrInvalidProtocolOrVersion
	}

	rl := &RequestLine{
		Method:        method,
		RequestTarget: reqTarget,
		HttpVersion:   protocol,
	}

	return rl, len(startLine), nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	bufLen := 0
	req := NewRequest()

	f, _ := os.Create("log.txt")
	defer f.Close()
	log.SetOutput(f)

	for !req.done() {
		if bufLen == cap(buf) {
			newBufPart := make([]byte, len(buf))
			buf = append(buf, newBufPart...)
		}
		readedBytes, err := reader.Read(buf[bufLen:])
		log.Printf("readed %d", readedBytes)

		if readedBytes == 0 && err == io.EOF {
			req.state = StateDone
			break
		}

		if err != nil {
			return nil, err
		}

		bufLen += readedBytes
		parsedBytes, err := req.parse(buf[:bufLen])

		if err != nil {
			return nil, err
		}

		copy(buf, buf[parsedBytes:bufLen])
		bufLen -= parsedBytes

	}
	return req, nil
}

func isUpperAndLetters(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) || !unicode.IsUpper(c) {
			return false
		}
	}
	return true
}
