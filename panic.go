package tigo

import (
	"runtime/debug"
	"io"
	"fmt"
	"time"
)

func Panic(writer io.Writer) Handler {
	return func(ctx *Context) error{
		defer func() {
			if err := recover(); err != nil{
				//abort next
				ctx.Abort()
				if writer == nil{
					return
				}
				writer.Write([]byte(fmt.Sprintf(
					"----------- tigo panic info start --------------\nError:%v\nTime:%v\nUri:%s\nRemote-Addr:%s\n%s\n%s----------- tigo panic info end --------------\n",
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