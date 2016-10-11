package tigo

import (
	"time"
	"fmt"
	"io"
	"os"
	"net/http"
	"strings"
)

func Logger(writer io.Writer) Handler {
	if writer == nil{
		writer = os.Stdout
	}
	return func(ctx *Context) error {
		start := time.Now()
		format := `{"time":"%v","method":"%s","uri":"%s","status":"%v","referer":"%s","host":"%s","user_agent":"%s","remote_addr":"%s","latency":"%s","request_length":"%v","response_length":"%v", "error":"%v"}` + "\n"
		err := ctx.Next()
		statusCode := ctx.Response.StatusCode()
		errmsg := ""
		if err != nil {
			if httpError, ok := err.(HTTPError); ok {
				statusCode = httpError.StatusCode()
			} else {
				statusCode = http.StatusInternalServerError
			}
			errmsg = err.Error()
		}
		latency := time.Now().Sub(start)
		writer.Write([]byte(fmt.Sprintf(format,
			time.Now().Format(time.RFC3339),
			ctx.Request.Header.Method(),
			ctx.Request.URI(),
			statusCode,
			ctx.Request.Header.Referer(),
			ctx.Request.Host(),
			ctx.Request.Header.UserAgent(),
			ctx.RemoteAddr(),
			latency,
			ctx.Request.Header.ContentLength(),
			len(ctx.Response.Body()),
			strings.Replace(errmsg, "\"", "\\\"", -1),
		)))

		return err
	}
}