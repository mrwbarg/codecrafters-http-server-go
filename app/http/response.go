package http

import (
	"fmt"
)

type Response struct {
	Version    float32
	StatusCode int
	Reason     string
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
