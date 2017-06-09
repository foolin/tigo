package tigo

import (
	"time"
	"fmt"
	"io"
	"os"
	"net/http"
	"strconv"
)

func Logger(writer io.Writer) Handler {
	if writer == nil{
		writer = os.Stdout
	}
	return func(ctx *Context) error {
		start := time.Now()

		rw := &LogResponseWriter{ctx.Response, http.StatusOK, 0}
		ctx.Response = rw

		format := `{"time":"%v","method":"%s","uri":"%s","status":"%v","referer":"%s","host":"%s","user_agent":%v,"remote_addr":"%s","latency":"%s","request_length":"%v","response_length":"%v", "error":%v}` + "\n"
		err := ctx.Next()
		statusCode := rw.Status
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
			ctx.Request.Method,
			ctx.Request.URL.RequestURI(),
			statusCode,
			ctx.Request.Referer(),
			ctx.Request.Host,
			strconv.Quote(fmt.Sprintf("%s", ctx.Request.UserAgent())),
			ctx.RequestIP(),
			latency,
			ctx.Request.ContentLength,
			rw.BytesWritten,
			strconv.Quote(errmsg),
		)))

		return err
	}
}


// LogResponseWriter wraps http.ResponseWriter in order to capture HTTP status and response length information.
type LogResponseWriter struct {
	http.ResponseWriter
	Status       int
	BytesWritten int64
}

func (r *LogResponseWriter) Write(p []byte) (int, error) {
	written, err := r.ResponseWriter.Write(p)
	r.BytesWritten += int64(written)
	return written, err
}

func (r *LogResponseWriter) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}