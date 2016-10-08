package tigo

import "github.com/valyala/fasthttp"

func (r *Router) Run(addr string) error {
	if addr == ""{
		addr = ":8080"
	}
	return fasthttp.ListenAndServe(addr, r.HandleRequest)
}