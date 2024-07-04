package main

import (
	"errors"
	"flag"
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
	conn.Write(res.WriteBytes(req.CompressResponse()))
	conn.Close()

}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	filesDir := flag.String("directory", "/", "Directory to serve files from")
	flag.Parse()

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
	router.Get("/files/:fileName", func(ctx *http.Context) *http.Response {
		res := &http.Response{}
		fileName, ok := ctx.PathArgs["fileName"].(string)
		if !ok {
			return res.WithVersion(1.1).WithStatusCode(400).WithReason("Bad Request")
		}
		fullDir := *filesDir + fileName

		data, err := os.ReadFile(fullDir)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return res.WithVersion(1.1).WithStatusCode(404).WithReason("Not Found")
			}
			return res.WithVersion(1.1).WithStatusCode(500).WithReason("Internal Server Error")

		}

		res = res.
			WithVersion(1.1).
			WithStatusCode(200).
			WithReason("OK").
			WithBody(string(data)).
			WithHeader("Content-Type", "application/octet-stream")
		return res
	})
	router.Post("/files/:fileName", func(ctx *http.Context) *http.Response {
		res := &http.Response{}
		fileName, ok := ctx.PathArgs["fileName"].(string)
		if !ok {
			return res.WithVersion(1.1).WithStatusCode(400).WithReason("Bad Request")
		}
		fullDir := *filesDir + fileName

		file, err := os.Create(fullDir)
		if err != nil {
			return res.WithVersion(1.1).WithStatusCode(500).WithReason("Internal Server Error")
		}
		defer file.Close()

		file.Write([]byte(ctx.Request.Body))
		file.Sync()

		res = res.
			WithVersion(1.1).
			WithStatusCode(201).
			WithReason("Created")
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
