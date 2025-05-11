package responsebuilder

import (
	"strconv"
)

func OK(body string) string {
	builder, _ := New(200)
	builder.AddHeader("Content-Type", "text/plain")
	builder.AddHeader("Content-Length", strconv.Itoa(len(body)))
	builder.SetBody(body)
	response, err := builder.Build()
	if err != nil {
		panic(err)
	}

	return response
}

func Created() string {
	builder, _ := New(201)
	response, err := builder.Build()
	if err != nil {
		panic(err)
	}

	return response
}

func NoContent() string {
	builder, _ := New(204)
	response, err := builder.Build()
	if err != nil {
		panic(err)
	}

	return response
}

func BadRequest() string {
	builder, _ := New(400)
	response, err := builder.Build()
	if err != nil {
		panic(err)
	}

	return response
}

func NotFound() string {
	builder, _ := New(404)
	response, err := builder.Build()
	if err != nil {
		panic(err)
	}

	return response
}

func InternalServerError() string {
	builder, _ := New(500)
	response, err := builder.Build()
	if err != nil {
		panic(err)
	}

	return response
}
