package main

import (
	"fmt"
	"net"
	"os"
)


type Response struct {
	Version float32 
	StatusCode int
	Reason string
}

func (o *Response) WithVersion(version float32) *Response {
	o.Version = version
	return o
}

func (o *Response) WithStatusCode(statusCode int) *Response {
	o.StatusCode = statusCode
	return o
}

func (o *Response) WithReason(reason string) *Response {
	o.Reason = reason
	return o
}

func (o *Response) WriteBytes() []byte {
	return []byte(fmt.Sprintf("HTTP/%.1f %d %s\r\n\r\n", o.Version, o.StatusCode, o.Reason))
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	r := &Response{}
	r = r.WithVersion(1.1).WithStatusCode(200).WithReason("OK")

	n, err := conn.Write(r.WriteBytes())

	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Wrote", n, "bytes")
}
