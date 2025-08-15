package request

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := min(cr.pos+cr.numBytesPerRead, len(cr.data))
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}
	return n, nil
}

func TestRequestLineParse(t *testing.T) {
	// t.Run("ok, GET Request line", func(t *testing.T) {
	// 	reader := &chunkReader{
	// 		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	// 		numBytesPerRead: 1,
	// 	}
	// 	r, err := RequestFromReader(reader)
	// 	require.NoError(t, err)
	// 	require.NotNil(t, r)
	// 	assert.Equal(t, "GET", r.RequestLine.Method)
	// 	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	// 	assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
	// })
	//
	// t.Run("ok, POST Request line", func(t *testing.T) {
	// 	reader := &chunkReader{
	// 		data:            "POST / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	// 		numBytesPerRead: 1,
	// 	}
	// 	r, err := RequestFromReader(reader)
	// 	require.NoError(t, err)
	// 	require.NotNil(t, r)
	// 	assert.Equal(t, "POST", r.RequestLine.Method)
	// 	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	// 	assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
	// })
	//
	// t.Run("ok, Request line with path", func(t *testing.T) {
	// 	reader := &chunkReader{
	// 		data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	// 		numBytesPerRead: 1,
	// 	}
	// 	r, err := RequestFromReader(reader)
	// 	require.NoError(t, err)
	// 	require.NotNil(t, r)
	// 	assert.Equal(t, "GET", r.RequestLine.Method)
	// 	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	// 	assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
	//
	// 	reader = &chunkReader{
	// 		data:            "GET /coffee?q=val HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	// 		numBytesPerRead: 1,
	// 	}
	// 	r, err = RequestFromReader(reader)
	// 	require.NoError(t, err)
	// 	require.NotNil(t, r)
	// 	assert.Equal(t, "GET", r.RequestLine.Method)
	// 	assert.Equal(t, "/coffee?q=val", r.RequestLine.RequestTarget)
	// 	assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
	//
	// 	reader = &chunkReader{
	// 		data:            "GET http://www.example.org/pub/WWW/TheProject.html HTTP/1.1\r\nHost: www.example.org\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	// 		numBytesPerRead: 1,
	// 	}
	// 	r, err = RequestFromReader(reader)
	// 	require.NoError(t, err)
	// 	require.NotNil(t, r)
	// 	assert.Equal(t, "GET", r.RequestLine.Method)
	// 	assert.Equal(t, "http://www.example.org/pub/WWW/TheProject.html", r.RequestLine.RequestTarget)
	// 	assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
	//
	// 	reader = &chunkReader{
	// 		data:            "GET www.example.com:80 HTTP/1.1\r\nHost: www.example.com:80\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	// 		numBytesPerRead: 1,
	// 	}
	// 	r, err = RequestFromReader(reader)
	// 	require.NoError(t, err)
	// 	require.NotNil(t, r)
	// 	assert.Equal(t, "GET", r.RequestLine.Method)
	// 	assert.Equal(t, "www.example.com:80", r.RequestLine.RequestTarget)
	// 	assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
	//
	// 	reader = &chunkReader{
	// 		data:            "OPTIONS * HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	// 		numBytesPerRead: 1,
	// 	}
	// 	r, err = RequestFromReader(reader)
	// 	require.NoError(t, err)
	// 	require.NotNil(t, r)
	// 	assert.Equal(t, "OPTIONS", r.RequestLine.Method)
	// 	assert.Equal(t, "*", r.RequestLine.RequestTarget)
	// 	assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
	// })

	t.Run("fail, Invalid number of parts in request line", func(t *testing.T) {
		reader := &chunkReader{
			data:            "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 1,
		}
		_, err := RequestFromReader(reader)
		require.Error(t, err)
		assert.ErrorIs(t, ErrMalformedRequestLine, err)

		reader = &chunkReader{
			data:            "GET /coffee HTTP/1.1 TEST \r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 1,
		}
		_, err = RequestFromReader(reader)
		require.Error(t, err)
		assert.ErrorIs(t, ErrMalformedRequestLine, err)
	})

	t.Run("fail, Invalid method(out of order) request line", func(t *testing.T) {
		reader := &chunkReader{
			data:            "gEt /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 1,
		}
		_, err := RequestFromReader(reader)
		require.Error(t, err)
		assert.ErrorIs(t, ErrInvalidMethodName, err)

		reader = &chunkReader{
			data:            "123 /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 1,
		}
		_, err = RequestFromReader(reader)
		require.Error(t, err)
		assert.ErrorIs(t, ErrInvalidMethodName, err)
	})

	t.Run("fail, Invalid version in request line", func(t *testing.T) {
		reader := &chunkReader{
			data:            "GET /coffee HTTP/2.0\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 1,
		}
		_, err := RequestFromReader(reader)
		require.Error(t, err)

		reader = &chunkReader{
			data:            "GET /coffee HTTP/3.0\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 1,
		}
		_, err = RequestFromReader(reader)
		require.Error(t, err)
		assert.ErrorIs(t, ErrInvalidProtocolOrVersion, err)

		reader = &chunkReader{
			data:            "GET /coffee FTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 1,
		}
		_, err = RequestFromReader(reader)
		require.Error(t, err)
		assert.ErrorIs(t, ErrInvalidProtocolOrVersion, err)
	})

	t.Run("fail, Invalid beginning of line", func(t *testing.T) {
		reader := &chunkReader{
			data:            "  GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 1,
		}
		_, err := RequestFromReader(reader)
		require.Error(t, err)
		assert.ErrorIs(t, ErrMalformedRequestLine, err)
	})
}
