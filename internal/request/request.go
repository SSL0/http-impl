package request

import (
	"errors"
	"io"
	"net/url"
	"strings"
	"unicode"
)

const separator = "\r\n"

var ErrMalformedRequestLine = errors.New("malformed request-line")
var ErrMissingSeparator = errors.New("missing separator")
var ErrInvalidMethodName = errors.New("incorrect method name")
var ErrIncorrectRequestTarget = errors.New("incorrect request target")
var ErrInvalidProtocolOrVersion = errors.New("incorrect protocol name or protocol version, only HTTP/1.1 allowed")

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func IsUpperAndLetters(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) || !unicode.IsUpper(c) {
			return false
		}
	}
	return true
}

func parseRequestLine(s string) (*RequestLine, string, error) {
	if s[0] == ' ' {
		return nil, s, ErrMalformedRequestLine
	}

	idx := strings.Index(s, separator)
	if idx == -1 {
		return nil, s, ErrMissingSeparator
	}

	startLine := s[:idx]
	restData := s[idx+len(separator):]

	startLineParts := strings.Split(startLine, " ")
	if len(startLineParts) != 3 {
		return nil, restData, ErrMalformedRequestLine
	}

	if !IsUpperAndLetters(startLineParts[0]) {
		return nil, restData, ErrInvalidMethodName
	}

	if _, err := url.Parse(startLineParts[1]); err != nil {
		return nil, restData, ErrIncorrectRequestTarget
	}

	protocolParts := strings.Split(startLineParts[2], "/")

	if len(protocolParts) != 2 || protocolParts[0] != "HTTP" || protocolParts[1] != "1.1" {
		return nil, restData, ErrInvalidProtocolOrVersion
	}

	rl := &RequestLine{
		Method:        startLineParts[0],
		RequestTarget: startLineParts[1],
		HttpVersion:   startLineParts[2],
	}

	return rl, restData, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	str := string(data)
	rl, _, err := parseRequestLine(str)
	if err != nil {
		return nil, err
	}

	return &Request{RequestLine: *rl}, nil
}
