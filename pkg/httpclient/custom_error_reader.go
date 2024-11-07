package httpclient

import (
	"errors"
	"io"
)

// CustomErrorReader is a custom io.Reader implementation that returns an error.
type CustomErrorReader struct {
	closed bool
}

func (c *CustomErrorReader) Read(p []byte) (n int, err error) {
	if c.closed {
		return 0, io.ErrClosedPipe
	}
	return 0, errors.New("simulated error during read")
}

func (c *CustomErrorReader) Close() error {
	c.closed = true
	return nil
}
