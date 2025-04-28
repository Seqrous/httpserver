package httprequest

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type HTTPMethod string

const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
)

type Request struct {
	Method  HTTPMethod
	Target  string
	Headers map[string]string
	Body    []byte
}

func Parse(r io.Reader) (*Request, error) {
	buffer := bufio.NewReader(r)
	line, err := buffer.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read request line: %v", err)
	}

	line = strings.TrimSpace(line)
	status := strings.Split(line, " ")
	if len(status) != 3 {
		return nil, fmt.Errorf("malformed status line: %v", line)
	}

	method, target := status[0], status[1]
	httpMethod, err := validateMethod(method)
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	for {
		line, err = buffer.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to parse headers: %v", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break // end of headers
		}

		header := strings.Split(line, ":")
		if len(header) != 2 {
			continue // ignore malformed headers
		}

		key := strings.TrimSpace(header[0])
		value := strings.TrimSpace(header[1])
		headers[key] = value
	}

	contentLength := 0
	if cl, ok := headers["Content-Length"]; ok {
		contentLength, err = strconv.Atoi(cl)
		if err != nil {
			return nil, fmt.Errorf("malformed Content-Length: %v", err)
		}
	}

	var body []byte
	if contentLength > 0 {
		body = make([]byte, contentLength)
		_, err = io.ReadFull(buffer, body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %v", err)
		}
	}

	return &Request{
		Method:  httpMethod,
		Target:  target,
		Headers: headers,
		Body:    body,
	}, nil
}

func validateMethod(method string) (HTTPMethod, error) {
	switch HTTPMethod(method) {
	case GET, POST, PUT, DELETE:
		return HTTPMethod(method), nil
	default:
		return "", errors.New("unsupported http method")
	}
}
