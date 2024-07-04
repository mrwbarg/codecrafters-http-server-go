package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

func handleConnection(conn net.Conn, router *http.Router) {
	fmt.Println("Handling new connection")

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("Read %d bytes.", n)
	req, err := http.ParseRequest(buffer)
	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		os.Exit(1)
	}

	res := router.Handle(req)
	conn.Write(res.WriteBytes())
	conn.Close()

}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	router := http.NewRouter()

	router.Get("/", func(ctx *http.Context) *http.Response {
		res := &http.Response{}
		res = res.WithVersion(1.1).WithStatusCode(200).WithReason("OK")
		return res
	})
	router.Get("/echo/:toEcho", func(ctx *http.Context) *http.Response {
		res := &http.Response{}
		toEcho, ok := ctx.PathArgs["toEcho"].(string)
		if !ok {
			return res.WithVersion(1.1).WithStatusCode(400).WithReason("Bad Request")
		}
		res = res.
			WithVersion(1.1).
			WithStatusCode(200).
			WithReason("OK").
			WithBody(toEcho).
			WithHeader("Content-Type", "text/plain")
		return res
	})
	router.Get("/user-agent", func(ctx *http.Context) *http.Response {
		res := &http.Response{}
		userAgent, ok := ctx.Request.Headers["User-Agent"]
		if !ok {
			return res.WithVersion(1.1).WithStatusCode(400).WithReason("Bad Request")
		}
		res = res.
			WithVersion(1.1).
			WithStatusCode(200).
			WithReason("OK").
			WithBody(userAgent).
			WithHeader("Content-Type", "text/plain")
		return res
	})

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn, router)
	}

}
