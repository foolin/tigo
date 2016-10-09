package tigo

import (
	"time"
	"fmt"
	"io"
	"os"
)

func Logger(writer io.Writer) Handler {
	if writer == nil{
		writer = os.Stdout
	}
	return func(ctx *Context) error {
		start := time.Now()
		format := `{"time":"%v","method":"%s","uri":"%s","status":"%v","referer":"%s","host":"%s","user_agent":"%s","remote_addr":"%s","latency":"%s","request_length":"%v","response_length":"%v"}` + "\n"
		defer func() {
			latency := time.Now().Sub(start)
			writer.Write([]byte(fmt.Sprintf(format,
				time.Now().Format(time.RFC3339),
				ctx.Request.Header.Method(),
				ctx.Request.URI(),
				ctx.Response.StatusCode(),
				ctx.Request.Header.Referer(),
				ctx.Request.Host(),
				ctx.Request.Header.UserAgent(),
				ctx.RemoteAddr(),
				latency,
				ctx.Request.Header.ContentLength(),
				len(ctx.Response.Body()),
			)))
		}()
		ctx.Next()
		return nil
	}
}