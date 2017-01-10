// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tigo

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	"errors"
	"strings"
	"net"
	"time"
	"net/http"
)
// Context represents the contextual data and environment while processing an incoming HTTP request.
type Context struct {
	*fasthttp.RequestCtx

	router   *Router
	pnames   []string               // list of route parameter names
	pvalues  []string               // list of parameter values corresponding to pnames
	data     map[string]interface{} // data items managed by Get and Set
	index    int                    // the index of the currently executing handler in handlers
	handlers []Handler              // the handlers associated with the current route
}

// Router returns the Router that is handling the incoming HTTP request.
func (c *Context) Router() *Router {
	return c.router
}

// Param returns the named parameter value that is found in the URL path matching the current route.
// If the named parameter cannot be found, an empty string will be returned.
func (c *Context) Param(name string) string {
	for i, n := range c.pnames {
		if n == name {
			return c.pvalues[i]
		}
	}
	return ""
}

// Get returns the named data item previously registered with the context by calling Set.
// If the named data item cannot be found, nil will be returned.
func (c *Context) Get(name string) interface{} {
	return c.data[name]
}

// Set stores the named data item in the context so that it can be retrieved later.
func (c *Context) Set(name string, value interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[name] = value
}

// Next calls the rest of the handlers associated with the current route.
// If any of these handlers returns an error, Next will return the error and skip the following handlers.
// Next is normally used when a handler needs to do some postprocessing after the rest of the handlers
// are executed.
func (c *Context) Next() error {
	c.index++
	for n := len(c.handlers); c.index < n; c.index++ {
		if err := c.handlers[c.index](c); err != nil {
			return err
		}
	}
	return nil
}

// Abort skips the rest of the handlers associated with the current route.
// Abort is normally used when a handler handles the request normally and wants to skip the rest of the handlers.
// If a handler wants to indicate an error condition, it should simply return the error without calling Abort.
func (c *Context) Abort() {
	c.index = len(c.handlers)
}

// URL creates a URL using the named route and the parameter values.
// The parameters should be given in the sequence of name1, value1, name2, value2, and so on.
// If a parameter in the route is not provided a value, the parameter token will remain in the resulting URL.
// Parameter values will be properly URL encoded.
// The method returns an empty string if the URL creation fails.
func (c *Context) URL(route string, pairs ...interface{}) string {
	if r := c.router.routes[route]; r != nil {
		return r.URL(pairs...)
	}
	return ""
}


// RequestHeader returns the request header's value
// accepts one parameter, the key of the header (string)
// returns string
func (c *Context) RequestHeader(key string) string {
	val := c.Request.Header.Peek(key)
	return string(val)
}

// RequestIP gets just the Remote Address from the client.
func (c *Context) RequestIP() string {
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.RequestCtx.RemoteAddr().String())); err == nil {
		return ip
	}
	return ""
}

// RemoteAddr is like RequestIP but it checks for proxy servers also, tries to get the real client's request IP
func (c *Context) RemoteAddr() string {
	header := string(c.RequestCtx.Request.Header.Peek("X-Real-Ip"))
	realIP := strings.TrimSpace(header)
	if realIP != "" {
		return realIP
	}
	realIP = string(c.RequestCtx.Request.Header.Peek("X-Forwarded-For"))
	idx := strings.IndexByte(realIP, ',')
	if idx >= 0 {
		realIP = realIP[0:idx]
	}
	realIP = strings.TrimSpace(realIP)
	if realIP != "" {
		return realIP
	}
	return c.RequestIP()
}

// QueryString returns the query value of a single key/name
func (c *Context) QueryString(key string) string {
	val := c.QueryArgs().Peek(key)
	return string(val)
}

// QueryMultiString returns query value associated with the given key.
func (c *Context) QueryMultiString(key string) []string {
	arrBytes := c.QueryArgs().PeekMulti(key)
	return arrBytes2Strs(arrBytes)
}

// PostMultiString returns the post data values as []string of a single key/name
func (c *Context) PostMultiString(name string) []string {
	arrBytes := c.PostArgs().PeekMulti(name)
	return arrBytes2Strs(arrBytes)
}

// PostString returns the post data value of a single key/name
// returns an empty string if nothing found
func (c *Context) PostString(name string) string {
	if v := c.PostMultiString(name); len(v) > 0 {
		return v[0]
	}
	return ""
}

// FormMultiValue returns form value associated with the given key.
//
// The value is searched in the following places:
//
//   * Query string.
//   * POST or PUT body.
//
// There are more fine-grained methods for obtaining form values:
//
//   * QueryArgs for obtaining values from query string.
//   * PostArgs for obtaining values from POST or PUT body.
//   * MultipartForm for obtaining values from multipart form.
//   * FormFile for obtaining uploaded files.
//
// The returned value is valid until returning from RequestHandler.
func (c *Context) FormMultiValue(key string) [][]byte {
	arr := make([][]byte, 0)
	mv := c.QueryArgs().PeekMulti(key)
	if len(mv) > 0 {
		arr = append(arr, mv...)
	}
	mv = c.PostArgs().PeekMulti(key)
	if len(mv) > 0 {
		arr = append(arr, mv...)
	}
	mf, err := c.MultipartForm()
	if err == nil && mf.Value != nil {
		mstrs := mf.Value[key]
		if len(mstrs) > 0 {
			for _, v := range mstrs {
				arr = append(arr, []byte(v))
			}
		}
	}
	return arr
}

// FormString returns a single value, as string, from post request's data
func (c *Context) FormString(name string) string {
	return string(c.FormValue(name))
}

// FormMultiString returns form value associated with the given key.
//
// The value is searched in the following places:
//
//   * Query string.
//   * POST or PUT body.
//
// There are more fine-grained methods for obtaining form values:
//
//   * QueryArgs for obtaining values from query string.
//   * PostArgs for obtaining values from POST or PUT body.
//   * MultipartForm for obtaining values from multipart form.
//   * FormFile for obtaining uploaded files.
//
// The returned value is valid until returning from RequestHandler.
func (c *Context) FormMultiString(key string) []string {
	arrBytes := c.FormMultiValue(key)
	return arrBytes2Strs(arrBytes)
}

// SetHeader set resposne header, use Add for setting multiple header values under the same key.
func (c *Context) SetHeader(key, value string) {
	c.Response.Header.Set(key, value)
}

// writeWithStatusCode writes the given data of arbitrary type to the response.
// The method calls the Serialize() method to convert the data into a byte array and then writes
// the byte array to the response.
func (c *Context) writeWithContentType(contentType string, bytes []byte) error {
	c.SetContentType(contentType)
	_, err := c.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

// JSON writes json values to the response.
func (c *Context) JSON(data interface{}) (err error) {
	var bytes []byte
	if bytes, err = json.Marshal(data); err == nil {
		err = c.writeWithContentType("application/json; charset=utf-8", bytes)
	}
	return
}

// Text writes text values to the response.
func (c *Context) Text(content string) error {
	return c.writeWithContentType("text/plain; charset=utf-8", []byte(content))
}

// HTML writes HTML values to the response.
func (c *Context) HTML(content string) error {
	return c.writeWithContentType("text/html; charset=utf-8", []byte(content))
}

func (c *Context) clientAllowsGzip() bool {
	if h := c.RequestHeader("Accept-Encoding"); h != "" {
		for _, v := range strings.Split(h, ";") {
			if strings.Contains(v, "gzip") {
				// we do Contains because sometimes browsers has the q=, we don't use it atm. || strings.Contains(v,"deflate"){
				return true
			}
		}
	}

	return false
}

// Gzip accepts bytes, which are compressed to gzip format and sent to the client
func (c *Context) Gzip(b []byte, status int) (err error) {
	c.RequestCtx.Response.Header.Add("Vary", "Accept-Encoding")
	if c.clientAllowsGzip() {
		_, err = fasthttp.WriteGzip(c.RequestCtx.Response.BodyWriter(), b)
		if err == nil {
			c.SetHeader("Content-Encoding", "gzip")
		}
	}
	return
}

// IsAjax returns true if this request is an 'ajax request'( XMLHttpRequest)
//
// Read more at: http://www.w3schools.com/ajax/
func (c *Context) IsAjax() bool {
	return c.RequestHeader("X-Requested-With") == "XMLHttpRequest"
}

// BindJson post for json
func (c *Context) BindJson(out interface{}) error {
	bytes := c.Request.Body()
	if len(bytes) <= 0 {
		return nil
	}
	return json.Unmarshal(bytes, &out)
}

// Render render with master
func (c *Context) Render(name string, data interface{}) error {
	return c.doRender(name, data, false)

}

// Render render only file
func (c *Context) RenderFile(name string, data interface{}) error {
	return c.doRender(name, data, true)
}

func (c *Context) doRender(name string, data interface{}, isRenderFile bool) error {
	if c.router.render == nil {
		return errors.New("Render engine not found.")
	}
	contentType := string(c.Response.Header.ContentType())
	if contentType == "" || !strings.Contains(contentType, "text/html"){
		c.SetContentType("text/html; charset=utf-8")
	}
	if isRenderFile{
		return c.router.render.RenderFile(c.Response.BodyWriter(), name, data)
	}else{
		return c.router.render.Render(c.Response.BodyWriter(), name, data)
	}
}

// init sets the request and response of the context and resets all other properties.
func (c *Context) init(ctx *fasthttp.RequestCtx) {
	c.RequestCtx = ctx
	c.data = nil
	c.index = -1
}

func arrBytes2Strs(arrBytes [][]byte) []string {
	arrStr := make([]string, len(arrBytes))
	for i, v := range arrBytes {
		arrStr[i] = string(v)
	}
	return arrStr
}

func (c *Context) GetCookieValue(key string) string {
	return string(c.Request.Header.Cookie(key))
}

func (c *Context) SetCookieValue(key string, value string, expire time.Time) {
	cookie := &fasthttp.Cookie{}
	cookie.SetKey(key)
	cookie.SetValue(value)
	cookie.SetExpire(expire)
	c.Response.Header.SetCookie(cookie)
}

func (c *Context) DelCookie(key string) {
	c.Response.Header.DelCookie(key)
}

func (c *Context) NotFound() error {
	return c.router.OnNotFound(c)
}

func (c *Context) Error(err error) {
	if httpError, ok := err.(HTTPError); ok {
		c.RequestCtx.Response.SetStatusCode(httpError.StatusCode())
	}else{
		c.RequestCtx.Response.SetStatusCode(http.StatusInternalServerError)
	}
	c.router.OnError(c, err)
	c.Abort()
}

func (ctx *Context) Redirect(uri string) {
	ctx.RequestCtx.Redirect(uri, http.StatusOK)
}