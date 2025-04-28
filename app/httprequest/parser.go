package httprequest

import (
	"errors"
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
	Body    string
}

func Parse(request string) (*Request, error) {
	splits := strings.Split(request, "\r\n")
	if len(splits) < 3 {
		return nil, errors.New("malformed request")
	}

	// parse status line
	statusSplits := strings.Split(splits[0], " ")
	method, err := validateMethod(statusSplits[0])
	if err != nil {
		return nil, err
	}
	target := statusSplits[1]

	// skip first and last splits (status and body)
	headers := make(map[string]string)
	for i := 1; i < len(splits)-2; i++ {
		header := strings.Split(splits[i], ":")
		if len(header) != 2 {
			continue // skip malformed headers
		}
		key := strings.TrimSpace(header[0])
		value := strings.TrimSpace(header[1])
		headers[key] = value
	}

	body := splits[len(splits)-1]
	return &Request{
		Method:  method,
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
