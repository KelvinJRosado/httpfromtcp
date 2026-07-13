package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("Host:         localhost:42069         \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 40, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	data = []byte("Host:         localhost:42069         \r\nFoo-Val:Bar\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 40, n)
	assert.False(t, done)
	// 2nd pass
	n2, done, err := headers.Parse(data[n:])
	require.NoError(t, err)
	assert.Equal(t, 13, n2)
	assert.False(t, done)
	// 3rd pass
	n3, done, err := headers.Parse(data[n+n2:])
	require.NoError(t, err)
	assert.Equal(t, 2, n3)
	assert.True(t, done)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "Bar", headers["foo-val"])

	// Test: Valid 2 headers with repeated headers
	headers = NewHeaders()
	data = []byte("FoO:Bar1\r\nfOo:Bar2\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 10, n)
	assert.False(t, done)
	// 2nd pass
	n2, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	assert.Equal(t, 10, n2)
	assert.False(t, done)
	// 3rd pass
	n3, done, err = headers.Parse(data[n+n2:])
	require.NoError(t, err)
	assert.Equal(t, 2, n3)
	assert.True(t, done)
	assert.Equal(t, "Bar1,Bar2", headers["foo"])

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Missing CRLF
	headers = NewHeaders()
	data = []byte("foo")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: No field name provided
	headers = NewHeaders()
	data = []byte(": localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: No ":" provided
	headers = NewHeaders()
	data = []byte("Hostlocalhost\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid char in header name
	headers = NewHeaders()
	data = []byte("H@st:localhost\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
