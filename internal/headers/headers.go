package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

const (
	CRLF                   = "\r\n"
	CRLFLen                = 1
	fieldValueInvalidChars = "\r\n\x00"
	fieldNameRegExp        = `^([\w-]+)$`
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

	value = strings.TrimSpace(value)

	re := regexp.MustCompile(fieldNameRegExp)
	// TODO: rewrite validation
	if !re.MatchString(key) {
		return "", "", fmt.Errorf("failed to validate field-name: %s", key)
	}

	if strings.ContainsAny(value, fieldValueInvalidChars) {
		return "", "", fmt.Errorf("value contains invalid chars: %s", value)
	}

	return key, value, nil

}
