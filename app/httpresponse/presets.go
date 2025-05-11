package httpresponse

import (
	"strconv"
)

func OK(body string) string {
	builder, _ := new(200)
	builder.addHeader("Content-Type", "text/plain")
	builder.addHeader("Content-Length", strconv.Itoa(len(body)))
	builder.setBody(body)
	response, err := builder.build()
	if err != nil {
		panic(err)
	}

	return response
}

func Created() string {
	builder, _ := new(201)
	response, err := builder.build()
	if err != nil {
		panic(err)
	}

	return response
}

func NoContent() string {
	builder, _ := new(204)
	response, err := builder.build()
	if err != nil {
		panic(err)
	}

	return response
}

func BadRequest() string {
	builder, _ := new(400)
	response, err := builder.build()
	if err != nil {
		panic(err)
	}

	return response
}

func NotFound() string {
	builder, _ := new(404)
	response, err := builder.build()
	if err != nil {
		panic(err)
	}

	return response
}

func InternalServerError() string {
	builder, _ := new(500)
	response, err := builder.build()
	if err != nil {
		panic(err)
	}

	return response
}
