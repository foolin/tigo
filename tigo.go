package tigo

import (
	"github.com/valyala/fasthttp"
	"os"
	"html/template"
)

// New creates a new Router object.
func New() *Router {
	r := &Router{
		routes: make(map[string]*Route),
		stores: make(map[string]routeStore),
	}
	r.RouteGroup = *newRouteGroup("", r, make([]Handler, 0))
	//r.NotFound(MethodNotAllowedHandler, NotFoundHandler)
	r.OnNotFound = NotFoundHandler
	r.OnError = HttpErrorHandler
	r.pool.New = func() interface{} {
		ctx := &Context{
			pvalues: make([]string, r.maxParams),
			router:  r,
		}
		return ctx
	}
	return r
}

//Default create default router, use panic, logger and render.
func Default() *Router {
	r := New()
	//logger
	r.Use(Logger(os.Stdout))
	//panic
	r.Use(Panic(os.Stderr))
	//tempate
	r.render = NewViewRender(ViewRenderConfig{
		Root: "views",
		Extension: ".html",
		Master: "master",
		Partials: []string{},
		Funcs: make(template.FuncMap),
		DisableCache: true,
	})
	return r
}

func (r *Router) Run(addr string) error {
	if addr == "" {
		addr = ":8080"
	}
	if r.render != nil {
		err := r.render.Init()
		if err != nil {
			return err
		}
	}
	return fasthttp.ListenAndServe(addr, r.HandleRequest)
}