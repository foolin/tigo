package tigo

import (
	"github.com/valyala/fasthttp"
)

func Static(root string, stripSlashes int) Handler {
	fsHandler := fasthttp.FSHandler(root, stripSlashes)
	return func(ctx *Context) error {
		//log.Printf("static file: %s", ctx.Request.URI().Path())
		fsHandler(ctx.RequestCtx)
		return nil
	}
}