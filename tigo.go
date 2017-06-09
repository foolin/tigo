package tigo

import (
	"os"
	"html/template"
	"net/http"
)

// ========================================== //
//	_____ _
//	|_   _(_) __ _  ___
//	| | | |/ _` |/ _ \
//	  | | | | (_| | (_) |
//	  |_| |_|\__, |\___/
//		 |___/
//
// ========================================== //

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
		DisableCache: false,
		DisableFilePartial: true,
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
	http.Handle("/", r)
	return http.ListenAndServe(addr, nil)
}