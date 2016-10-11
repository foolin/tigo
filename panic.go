package tigo

import (
	"runtime/debug"
	"io"
	"fmt"
	"time"
	"net/http"
)

func Panic(writer io.Writer) Handler {
	return func(ctx *Context) (err error){
		defer func() {
			if rerr := recover(); rerr != nil{
				err = NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("%s", rerr))
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
				//abort next
				ctx.Abort()
			}
		}()
		err = ctx.Next()
		return
	}
}