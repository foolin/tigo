package tigo

import (
	"runtime/debug"
	"io"
	"fmt"
	"time"
	"net/http"
)

func Panic(writer io.Writer) Handler {
	return func(ctx *Context) error{
		defer func() {
			if err := recover(); err != nil{
				//abort next
				ctx.Abort()
				ctx.SetStatusCode(http.StatusInternalServerError)
				if writer == nil{
					return
				}
				writer.Write([]byte(fmt.Sprintf(
					"----------- Tigo panic info start --------------\nError:%v\nTime:%v\nUri:%s\nRemote-Addr:%s\n%s\n%s----------- Tigo panic info end --------------\n",
					err,
					time.Now().Format(time.RFC3339),
					ctx.Request.URI(),
					ctx.RemoteAddr(),
					ctx.Request.Header.Header(),
					debug.Stack(),
				)))
			}
		}()
		ctx.Next()
		return nil
	}
}