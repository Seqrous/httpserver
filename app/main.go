package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"httpserver/app/httprequest"
	"httpserver/app/responsebuilder"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleRequestWithRecovery(conn)
	}
}

func handleRequestWithRecovery(conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			builder, _ := responsebuilder.New(500)
			response, _ := builder.Build()
			conn.Write([]byte(response))
		}
		conn.Close()
	}()
	handleRequest(conn)
}

func handleRequest(conn net.Conn) {
	request, err := httprequest.Parse(conn)
	if err != nil {
		builder, _ := responsebuilder.New(400)
		response, _ := builder.Build()
		conn.Write([]byte(response))
		return
	}

	if request.Target == "/" {
		builder, _ := responsebuilder.New(200)
		response, _ := builder.Build()
		conn.Write([]byte(response))
	} else if strings.HasPrefix(request.Target, "/user-agent") {
		responseBody := request.Headers["User-Agent"]
		fmt.Println(responseBody)
		builder, _ := responsebuilder.New(200)
		builder.AddHeader("Content-Type", "text/plain")
		builder.AddHeader("Content-Length", strconv.Itoa(len(responseBody)))
		builder.SetBody(responseBody)
		response, _ := builder.Build()
		conn.Write([]byte(response))
	} else if strings.HasPrefix(request.Target, "/echo/") {
		responseBody := strings.Split(request.Target, "/echo/")[1]
		builder, _ := responsebuilder.New(200)
		builder.AddHeader("Content-Type", "text/plain")
		builder.AddHeader("Content-Length", strconv.Itoa(len(responseBody)))
		builder.SetBody(responseBody)
		response, _ := builder.Build()
		conn.Write([]byte(response))
	} else if strings.HasPrefix(request.Target, "/files/") {
		fileName := strings.Split(request.Target, "/files/")[1]
		fileBuffer, err := os.ReadFile(fileName)
		if err != nil {
			builder, _ := responsebuilder.New(404)
			response, _ := builder.Build()
			conn.Write([]byte(response))
		}
		fileContent := string(fileBuffer)
		builder, _ := responsebuilder.New(200)
		builder.AddHeader("Content-Type", "text/html; charset=utf-8")
		builder.AddHeader("Content-Length", strconv.Itoa(len(fileContent)))
		builder.SetBody(fileContent)
		response, _ := builder.Build()
		conn.Write([]byte(response))
	} else {
		builder, _ := responsebuilder.New(404)
		response, _ := builder.Build()
		conn.Write([]byte(response))
	}
}
