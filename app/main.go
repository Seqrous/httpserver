package main

import (
	"fmt"
	"net"
	"os"
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
			response := responsebuilder.InternalServerError()
			conn.Write([]byte(response))
		}
		conn.Close()
	}()
	handleRequest(conn)
}

func handleRequest(conn net.Conn) {
	request, err := httprequest.Parse(conn)
	if err != nil {
		response := responsebuilder.BadRequest()
		conn.Write([]byte(response))
		return
	}

	if request.Target == "/" {
		response := responsebuilder.NoContent()
		conn.Write([]byte(response))
	} else if strings.HasPrefix(request.Target, "/user-agent") {
		responseBody := request.Headers["User-Agent"]
		fmt.Println(responseBody)
		response := responsebuilder.OK(responseBody)
		conn.Write([]byte(response))
	} else if strings.HasPrefix(request.Target, "/echo/") {
		responseBody := strings.Split(request.Target, "/echo/")[1]
		response := responsebuilder.OK(responseBody)
		conn.Write([]byte(response))
	} else if strings.HasPrefix(request.Target, "/files/") && request.Method == httprequest.GET {
		fileName := strings.Split(request.Target, "/files/")[1]
		fileBuffer, err := os.ReadFile(fileName)
		if err != nil {
			response := responsebuilder.NotFound()
			conn.Write([]byte(response))
		}
		fileContent := string(fileBuffer)
		response := responsebuilder.OK(fileContent)
		conn.Write([]byte(response))
	} else if strings.HasPrefix(request.Target, "/files/") && request.Method == httprequest.POST {
		fileName := strings.Split(request.Target, "/files/")[1]
		file, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}
		_, err = file.Write(request.Body)
		if err != nil {
			panic(err)
		}
		response := responsebuilder.Created()
		conn.Write([]byte(response))
	} else {
		response := responsebuilder.NotFound()
		conn.Write([]byte(response))
	}
}
