package tigo

import (
	"github.com/valyala/fasthttp"
	"os"
)

// New creates a new Router object.
func New() *Router {
	r := &Router{
		routes: make(map[string]*Route),
		stores: make(map[string]routeStore),
	}
	r.RouteGroup = *newRouteGroup("", r, make([]Handler, 0))
	r.NotFound(MethodNotAllowedHandler, NotFoundHandler)
	r.pool.New = func() interface{} {
		return &Context{
			pvalues: make([]string, r.maxParams),
			router:  r,
		}
	}
	return r
}

func Default() *Router {
	r := New()
	//panic
	r.Use(Panic(os.Stderr))
	//logger
	r.Use(Logger(os.Stdout))
	return r
}

func (r *Router) Run(addr string) error {
	if addr == ""{
		addr = ":8080"
	}
	return fasthttp.ListenAndServe(addr, r.HandleRequest)
}