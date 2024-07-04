package http

type Context struct {
	Request  *Request
	PathArgs map[string]any
}
