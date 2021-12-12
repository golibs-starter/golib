package client

import (
	"io"
)

type RequestReader interface {

	// Read the request
	Read() (io.Reader, error)
}
