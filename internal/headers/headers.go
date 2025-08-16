package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	CRLF                   = "\r\n"
	CRLFLen                = 1
	fieldValueInvalidChars = "\r\n\x00"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if len(data) == 0 {
		return 0, false, nil
	}
	sep := []byte(CRLF)
	sepLen := len(sep)

	dataIdx := 0

	for {
		sepIdx := bytes.Index(data, sep)

		if sepIdx == -1 {
			break
		}
		if sepIdx == 0 {
			n += CRLFLen
			done = true
			break
		}

		line := data[dataIdx:sepIdx]

		key, value, err := parseFieldLine(line)

		if err != nil {
			return n, done, fmt.Errorf("failed to parse field-line: %v", err)
		}

		if h[key] != "" {
			h[key] += ", "
		}
		h[key] += value

		data = data[sepIdx+sepLen:]
		n += len(line) + CRLFLen
	}

	return n, done, nil
}

func parseFieldLine(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte{':'}, 2)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("failed to split in to parts line: %s, parts %d", fieldLine, len(parts))
	}

	key := string(parts[0])
	value := string(parts[1])

	if !isToken(key) {
		return "", "", fmt.Errorf("failed to validate field-name: %s", key)
	}

	key = strings.ToLower(key)
	value = strings.TrimSpace(value)

	if strings.ContainsAny(value, fieldValueInvalidChars) {
		return "", "", fmt.Errorf("value contains invalid chars: %s", value)
	}

	return key, value, nil
}

func isToken(s string) bool {
	for _, ch := range s {
		valid := false
		if (ch >= 'A' && ch <= 'Z') ||
			(ch >= 'a' && ch <= 'z') ||
			(ch >= '0' && ch <= '9') {
			valid = true
		}

		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			valid = true
		}
		if !valid {
			fmt.Printf("VALIDATE ERROR: %s, %d", string(ch), ch)
			return false
		}

	}
	return true

}
