package http

import (
	"regexp"
	"sort"
	"strings"
)

type Route struct {
	RawPath string
	Handler func(*Context) *Response
}

func (r *Route) GetPathArgs(target string) map[string]any {
	pathNodes := strings.Split(r.RawPath, "/")
	targetNodes := strings.Split(target, "/")

	args := make(map[string]any)
	for i, node := range pathNodes {
		if strings.HasPrefix(node, ":") {
			args[node[1:]] = targetNodes[i]
		}
	}

	return args
}

func (r *Route) GetContext(req *Request) *Context {
	return &Context{
		Request:  req,
		PathArgs: r.GetPathArgs(req.Target),
	}
}

func (r *Route) Handle(ctx *Context) *Response {
	return r.Handler(ctx)
}

type Router struct {
	GET map[string]Route
}

func NewRouter() *Router {
	router := &Router{}
	router.GET = make(map[string]Route)
	return router
}

func pathToRegex(path string) string {
	nodes := strings.Split(path, "/")

	formattedNodes := make([]string, 0, len(nodes))
	for _, node := range nodes {
		if strings.HasPrefix(node, ":") {
			node = "(.*)"
		}
		formattedNodes = append(formattedNodes, node)
	}

	return strings.Join(formattedNodes, "/")
}

func (r *Router) Get(path string, handler func(*Context) *Response) {
	route := Route{
		RawPath: path,
		Handler: handler,
	}
	r.GET[pathToRegex(path)] = route
}

func (r *Router) Match(req *Request) (Route, *Context, bool) {
	paths := make([]string, 0, len(r.GET))

	for k := range r.GET {
		paths = append(paths, k)
	}

	sort.Strings(paths)

	for _, regexPath := range paths {
		match, _ := regexp.MatchString("^"+regexPath+"$", req.Target)
		if match {
			route := r.GET[regexPath]
			return route, route.GetContext(req), true
		}
	}
	return Route{}, nil, false
}

func (r *Router) Handle(req *Request) *Response {

	res := &Response{}
	if req.Method == GET {
		handler, context, ok := r.Match(req)
		if ok {
			return handler.Handle(context)
		}
		return res.WithVersion(1.1).WithStatusCode(404).WithReason("Not Found")
	}
	return res.WithVersion(1.1).WithStatusCode(404).WithReason("Not Found")
}
