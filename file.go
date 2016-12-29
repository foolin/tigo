package tigo

import (
	"github.com/valyala/fasthttp"
	"log"
)

func Static(root string, stripSlashes int) Handler {
	fsHandler := fasthttp.FSHandler(root, stripSlashes)
	return func(ctx *Context) error {
		log.Printf("static file: %s", ctx.Request.URI().Path())
		fsHandler(ctx.RequestCtx)
		return nil
	}
}


func File(path string) Handler {
	return func(ctx *Context) error {
		fasthttp.ServeFile(ctx.RequestCtx, path)
		return nil
	}
}