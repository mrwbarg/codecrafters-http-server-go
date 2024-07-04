package http

import (
	"fmt"
	"strconv"
	"strings"
)

type Request struct {
	Method  string
	Target  string
	Headers map[string]string
	Body    string
}

func ParseRequest(buffer []byte) (*Request, error) {
	strBuffer := string(buffer)

	requestElements := strings.Split(strBuffer, "\r\n")

	if len(requestElements) < 3 {
		return nil, fmt.Errorf("invalid request")
	}

	req := &Request{}

	requestData := strings.Split((requestElements[0]), " ")
	req.Method = requestData[0]
	req.Target = requestData[1]

	req.Headers = make(map[string]string)
	next := 1
	for {
		header := requestElements[next]
		if header == "" {
			break
		}

		headerElements := strings.SplitN(header, ":", 2)
		req.Headers[headerElements[0]] = strings.TrimSpace(headerElements[1])

		next++
	}
	next++

	if req.Method == POST {
		contentLength, ok := req.Headers["Content-Length"]
		rawBody := requestElements[next]
		if !ok && (len(rawBody) > 0) {
			return nil, fmt.Errorf("invalid request")
		} else {
			i, err := strconv.Atoi(contentLength)
			if err != nil {
				return nil, fmt.Errorf("invalid request")
			}
			req.Body = string([]byte(rawBody)[0:i])
		}
	}

	return req, nil
}
