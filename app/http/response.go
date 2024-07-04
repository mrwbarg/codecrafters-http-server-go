package http

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"strings"
)

type Response struct {
	Version    float32
	StatusCode int
	Reason     string
	Headers    map[string]any
	Body       string
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

func (o *Response) WithHeader(key string, value any) *Response {
	if o.Headers == nil {
		o.Headers = make(map[string]any)
	}
	o.Headers[key] = value
	return o
}

func (o *Response) SetHeader(key string, value any) *Response {
	if o.Headers == nil {
		o.Headers = make(map[string]any)
	}
	o.Headers[key] = value
	return o
}

func (o *Response) WithBody(body string) *Response {
	o.Body = body
	return o
}

func (o *Response) WriteHeaders() string {
	headers := make([]string, 0)
	for key, value := range o.Headers {
		headers = append(headers, fmt.Sprintf("%s: %v\r\n", key, value))
	}
	return strings.Join(headers, "")
}

func (o *Response) WriteBody(compress bool) string {
	if compress {
		var buffer bytes.Buffer
		gz := gzip.NewWriter(&buffer)
		if _, err := gz.Write([]byte(o.Body)); err != nil {
			return ""
		}
		if err := gz.Close(); err != nil {
			return ""
		}
		o.SetHeader("Content-Encoding", "gzip")
		return buffer.String()
	}
	return o.Body
}

func (o *Response) WriteBytes(compress bool) []byte {
	body := o.WriteBody(compress)
	o.SetHeader("Content-Length", len([]byte(body)))
	return []byte(
		fmt.Sprintf(
			"HTTP/%.1f %d %s\r\n%s\r\n%s",
			o.Version,
			o.StatusCode,
			o.Reason,
			o.WriteHeaders(),
			body,
		),
	)
}
