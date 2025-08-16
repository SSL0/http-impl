package headers

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const (
	CRLF                   = "\r\n"
	CRLFLen                = 2
	fieldValueInvalidChars = "\r\n\x00"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h *Headers) GetString(key string) (string, bool) {
	v, ok := (*h)[strings.ToLower(key)]
	return v, ok
}

func (h *Headers) GetInt(key string, defaultValue int) int {
	v, ok := (*h)[strings.ToLower(key)]

	if !ok {
		return defaultValue
	}

	i, err := strconv.Atoi(v)

	if err != nil {
		return defaultValue
	}

	return i
}

func (h *Headers) Set(key, value string) {
	key = strings.ToLower(key)

	if v, ok := (*h)[key]; ok {
		(*h)[key] = fmt.Sprintf("%s, %s", v, value)
	} else {
		(*h)[key] = value
	}
}

func (h *Headers) ForEach(callback func(k, v string)) {
	for k, v := range *h {
		callback(k, v)
	}
}

func (h *Headers) Parse(data []byte) (n int, done bool, err error) {
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

		h.Set(key, value)

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
			return false
		}

	}
	return true

}
