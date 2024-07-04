package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

func handleConnection(conn net.Conn, router *Router) {
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

type Router struct {
	GET map[string]func(*http.Request) *http.Response
}

func NewRouter() *Router {
	router := &Router{}
	router.GET = make(map[string]func(*http.Request) *http.Response)
	return router
}

func (r *Router) get(path string, handler func(*http.Request) *http.Response) {
	r.GET[path] = handler
}

func (r *Router) Handle(req *http.Request) *http.Response {

	res := &http.Response{}
	if req.Method == http.GET {
		handler, ok := r.GET[req.Target]
		if ok {
			return handler(req)
		}
		return res.WithVersion(1.1).WithStatusCode(404).WithReason("Not Found")
	}
	return res.WithVersion(1.1).WithStatusCode(404).WithReason("Not Found")
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	router := NewRouter()
	router.get("/", func(req *http.Request) *http.Response {
		res := &http.Response{}
		res = res.WithVersion(1.1).WithStatusCode(200).WithReason("OK")
		return res
	})

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		handleConnection(conn, router)
	}

}
