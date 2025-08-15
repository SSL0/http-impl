package request

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse(t *testing.T) {
	t.Run("ok, GET Request line", func(t *testing.T) {
		r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/", r.RequestLine.RequestTarget)
		assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
	})
	t.Run("ok, POST Request line", func(t *testing.T) {
		r, err := RequestFromReader(strings.NewReader("POST / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "POST", r.RequestLine.Method)
		assert.Equal(t, "/", r.RequestLine.RequestTarget)
		assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
	})
	t.Run("ok, Request line with path", func(t *testing.T) {
		r, err := RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
		assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)

		r, err = RequestFromReader(strings.NewReader("GET /coffee?q=val HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/coffee?q=val", r.RequestLine.RequestTarget)
		assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)

		r, err = RequestFromReader(strings.NewReader("GET http://www.example.org/pub/WWW/TheProject.html HTTP/1.1\r\nHost: www.example.org\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "http://www.example.org/pub/WWW/TheProject.html", r.RequestLine.RequestTarget)
		assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)

		r, err = RequestFromReader(strings.NewReader("GET www.example.com:80 HTTP/1.1\r\nHost: www.example.com:80\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "www.example.com:80", r.RequestLine.RequestTarget)
		assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)

		r, err = RequestFromReader(strings.NewReader("OPTIONS * HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "OPTIONS", r.RequestLine.Method)
		assert.Equal(t, "*", r.RequestLine.RequestTarget)
		assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
	})
	t.Run("fail, Invalid number of parts in request line", func(t *testing.T) {
		_, err := RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.Error(t, err)
		assert.ErrorIs(t, ErrMalformedRequestLine, err)

		_, err = RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1 TEST \r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.Error(t, err)
		assert.ErrorIs(t, ErrMalformedRequestLine, err)
	})
	t.Run("fail, Invalid method(out of order) request line", func(t *testing.T) {
		_, err := RequestFromReader(strings.NewReader("gEt /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.Error(t, err)
		assert.ErrorIs(t, ErrInvalidMethodName, err)

		_, err = RequestFromReader(strings.NewReader("123 /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.Error(t, err)
		assert.ErrorIs(t, ErrInvalidMethodName, err)
	})
	t.Run("fail, Invalid version in request line", func(t *testing.T) {
		_, err := RequestFromReader(strings.NewReader("GET /coffee HTTP/2.0\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.Error(t, err)

		_, err = RequestFromReader(strings.NewReader("GET /coffee HTTP/3.0\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.Error(t, err)
		assert.ErrorIs(t, ErrInvalidProtocolOrVersion, err)

		_, err = RequestFromReader(strings.NewReader("GET /coffee FTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.Error(t, err)
		assert.ErrorIs(t, ErrInvalidProtocolOrVersion, err)
	})
	t.Run("fail, Invalid beginning of line", func(t *testing.T) {
		_, err := RequestFromReader(strings.NewReader("  GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.Error(t, err)
		assert.ErrorIs(t, ErrMalformedRequestLine, err)
	})
}
