package request

import "io"

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	_, err := io.ReadAll(reader)

	var req *Request

	if err != nil {
		return req, err
	}

	return req, nil
}
