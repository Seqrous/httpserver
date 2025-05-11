package httpresponse

import (
	"fmt"
	"strings"
)

type responseBuilder struct {
	status  string
	headers map[string]string
	body    string
	err     error
}

func new(code int) (*responseBuilder, error) {
	var status string
	switch code {
	case 200:
		status = "OK"
	case 201:
		status = "Created"
	case 204:
		status = "No Content"
	case 400:
		status = "Bad Request"
	case 404:
		status = "Not Found"
	case 500:
		status = "Internal Server Error"
	default:
		return nil, fmt.Errorf("unsupported status code: %d", code)
	}

	return &responseBuilder{
		status:  fmt.Sprintf("HTTP/1.1 %d %s", code, status),
		headers: make(map[string]string),
	}, nil
}

func (rb *responseBuilder) addHeader(name string, value string) *responseBuilder {
	if rb.err != nil {
		return rb
	}

	rb.headers[name] = value
	return rb
}

func (rb *responseBuilder) setBody(body string) *responseBuilder {
	if rb.err != nil {
		return rb
	}

	if body == "" {
		rb.err = fmt.Errorf("body cannnot be empty")
		return rb
	}
	rb.body = body
	return rb
}

func (rb *responseBuilder) build() (string, error) {
	// return any error from previous steps
	if rb.err != nil {
		return "", rb.err
	}

	var b strings.Builder
	b.WriteString(rb.status)
	b.WriteString("\r\n")
	for key, value := range rb.headers {
		b.WriteString(key)
		b.WriteString(":")
		b.WriteString(value)
		b.WriteString("\r\n")
	}
	b.WriteString("\r\n")
	b.WriteString(rb.body)

	return b.String(), nil
}
