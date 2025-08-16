package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	t.Run("ok, single header", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost:42069\r\n\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "localhost:42069", headers["Host"])
		assert.Equal(t, len(data)-2, n)
		assert.True(t, done)
	})
	t.Run("ok, mutiple header", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost\r\nUser-Agent: curl/7.81.0\r\n\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "localhost", headers["Host"])
		assert.Equal(t, "curl/7.81.0", headers["User-Agent"])
		assert.Equal(t, len(data)-3, n)
		assert.True(t, done)
	})
	t.Run("ok, single header with extra whitespaces", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host:      localhost:42069     \r\n\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "localhost:42069", headers["Host"])
		assert.Equal(t, len(data)-2, n)
		assert.True(t, done)
	})
	t.Run("ok, two headers with existing headers", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Example-Field: Foo, Bar\r\nExample-Field: Baz\r\n\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "Foo, Bar, Baz", headers["Example-Field"])
		assert.Equal(t, len(data)-3, n)
		assert.True(t, done)
	})
	t.Run("ok, done false", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost:42069\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "localhost:42069", headers["Host"])
		assert.Equal(t, len(data)-1, n)
		assert.False(t, done)
	})
	t.Run("fail, spacing header", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("       Host : localhost:42069       \r\n\r\n")
		n, done, err := headers.Parse(data)
		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
	})
}
