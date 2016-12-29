package tigo

import (
	"github.com/valyala/fasthttp"
	"net/http"
)

func Static(root string, stripSlashes int) Handler {
	fsHandler := fasthttp.FSHandler(root, stripSlashes)
	return func(ctx *Context) error {
		fsHandler(ctx.RequestCtx)
		if ctx.RequestCtx.Response.StatusCode() == http.StatusNotFound {
			ctx.NotFound()
		}
		return nil
	}
}


func File(path string) Handler {
	return func(ctx *Context) error {
		fasthttp.ServeFile(ctx.RequestCtx, path)
		return nil
	}
}