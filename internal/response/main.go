package response

import "io"

type StatusCode int

const (
	Status200 StatusCode = iota
	Status400
	Status500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error
