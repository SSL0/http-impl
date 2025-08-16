package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"github.com/SSL0/http-impl/internal/headers"
)

const (
	bufferSize = 1024
	CRLF       = "\r\n"
	CRLFLen    = 2
)

type parserState string

const (
	stateRequestLine parserState = "init"
	stateHeaders     parserState = "headers"
	stateBody        parserState = "body"
	stateDone        parserState = "done"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       parserState
}

func NewRequest() *Request {
	return &Request{
		Headers: headers.NewHeaders(),
		Body:    []byte{},
		state:   stateRequestLine}
}

func (r *Request) done() bool {
	return r.state == stateDone
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	bufLen := 0
	req := NewRequest()

	for !req.done() {
		if bufLen == cap(buf) {
			newBufPart := make([]byte, len(buf))
			buf = append(buf, newBufPart...)
		}
		readedBytes, err := reader.Read(buf[bufLen:])

		if readedBytes == 0 && err == io.EOF {
			req.state = stateDone
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

	if req.Headers.GetInt("Content-Length", 0) != len(req.Body) {
		return nil, fmt.Errorf("body not equal content-length")
	}

	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalParsedBytes := 0
outer:
	for {
		currentData := data[totalParsedBytes:]
		if len(currentData) == 0 {
			break
		}
		switch r.state {
		case stateRequestLine:
			rl, n, err := parseRequestLine(currentData)

			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			totalParsedBytes += n

			r.state = stateHeaders
		case stateHeaders:
			n, done, err := r.Headers.Parse(currentData)

			if err != nil {
				return totalParsedBytes, err
			}

			if n == 0 {
				break outer
			}
			totalParsedBytes += n

			if done {
				r.state = stateBody
			}

		case stateBody:
			contentLengthValue, ok := r.Headers.GetString("Content-Length")
			if !ok || contentLengthValue == "" {
				r.state = stateDone
				break outer
			}

			contentLength, err := strconv.Atoi(contentLengthValue)
			if err != nil || contentLength == 0 {
				r.state = stateDone
				break outer
			}

			remaining := min(contentLength-len(r.Body), len(currentData))
			r.Body = append(r.Body, currentData[:remaining]...)
			totalParsedBytes += remaining

			if len(r.Body) == contentLength {
				r.state = stateDone
			}
		case stateDone:
			break outer
		default:
			return 0, errors.New("unknown parser state")
		}
	}
	return totalParsedBytes, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	if len(data) == 0 {
		return nil, 0, nil
	}

	if data[0] == ' ' {
		return nil, 0, fmt.Errorf("data starts from whitespace")
	}

	idx := bytes.Index(data, []byte(CRLF))
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := data[:idx]

	startLineParts := bytes.Split(startLine, []byte{' '})
	if len(startLineParts) != 3 {
		return nil, 0, fmt.Errorf("start line parts not equal three: %s", startLine)
	}

	method := string(startLineParts[0])
	reqTarget := string(startLineParts[1])
	protocol := string(startLineParts[2])

	if !isUpperAndLetters(method) {
		return nil, 0, fmt.Errorf("method contains unsupported chars: %s", method)
	}

	if _, err := url.Parse(reqTarget); err != nil {
		return nil, 0, fmt.Errorf("failed to parse url from request target: %s", reqTarget)
	}

	protocolParts := strings.Split(protocol, "/")

	if len(protocolParts) != 2 || protocolParts[0] != "HTTP" || protocolParts[1] != "1.1" {
		return nil, 0, fmt.Errorf("invalid protocol or protocol version: %s", protocol)
	}

	rl := &RequestLine{
		Method:        method,
		RequestTarget: reqTarget,
		HttpVersion:   protocol,
	}

	return rl, len(startLine) + CRLFLen, nil
}

func isUpperAndLetters(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) || !unicode.IsUpper(c) {
			return false
		}
	}
	return true
}
