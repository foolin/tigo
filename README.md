# Tigo


[![GoDoc](https://godoc.org/github.com/foolin/tigo?status.png)](https://godoc.org/github.com/foolin/tigo)

## Description

Tigo is a tiny framework for go, is forked from [Fasthttp-routing](https://github.com/qiangxue/fasthttp-routing) which is a lightweight high performance HTTP request router, fast and powerful routing features for the high-performance.
The package has the following features:

* middleware pipeline architecture, similar to that of the [Express framework](http://expressjs.com).
* extremely fast request routing with zero dynamic memory allocation
* modular code organization through route grouping
* flexible URL path matching, supporting URL parameters and regular expressions
* URL creation according to the predefined routes

This fork adds the following features:
* add OnError handdler for panic.
* add render function for view(use html/template)
* add panic middleware
* add logger middleware
* add static middleware
* add more router functions.

## Requirements

Go 1.5 or above.

## Installation

Run the following command to install the package:

```
go get github.com/foolin/tigo
```


## Getting Started

#### Default example:
```go
	//new default router
	router := tigo.Default()

	//register router
	router.Get("/", func(ctx *tigo.Context) error {
		content := `
			Hello tigo!!!<hr>
			visit api: <a href="/api/done">api/done</a>
		`
		//out html
		return ctx.HTML(content)
	})

	router.Get("/api/<action>", func(ctx *tigo.Context) error {

		//json object
		data := struct {
			Ip string `json:"ip"`
			Action string `json:"action"`
		}{ctx.RequestIP(), ctx.Param("action")}

		//out json
		return ctx.JSON(data)
	})

	//run
	log.Printf("run on :8080")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("run error: %v", err)
	}
```


#### New example:
```go

	//new router
	router := tigo.New()

	//logger
	router.Use(tigo.Logger(os.Stdout))

	//panic
	router.Use(tigo.Panic(os.Stderr))

	//register router
	router.Get("/", func(ctx *tigo.Context) error {
		return ctx.HTML("Hello tigo!!!")
	})

	//run
	log.Printf("run on :8080")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("run error: %v", err)
	}
```


#### Render example:

Use context.Render()

```go

	//new router
	router := tigo.New()

	//set render, tigo.Default() will default initialize.
	router.SetRender(tigo.NewHtmlRender(tigo.HtmlRenderConfig{
		ViewRoot:  "views",
		Extension: ".html",
	}))

	//register router
	router.Get("/", func(ctx *tigo.Context) error {
		return ctx.Render("page", tigo.M{"title": "Tigo render"})
	})

	//run
	log.Printf("run on :8080")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("run error: %v", err)
	}
	
```

/views/page.html

```html
		    <!-- /views/page.html content -->

		    <!doctype html>

		    <html>
		    <head>
			<meta http-equiv="Content-type" content="text/html; charset=utf-8" />
			<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
			<title>{{.title}}</title>
		    </head>

		    <body>
			page.html
			<hr>
			{{render "layout/footer"}}
		    </body>
		    </html>
```

/views/layout/footer.html
```html
		    <!-- /views/layout/footer.html content -->

		    Copyright &copy2016 by <a href="https://github.com/foolin/tigo">tigo</a>.
```


#### Render master example:

Use render with master page

```go

	//new router
    router := tigo.Default()

    admin := router.Group("/admin")

    //register router
    router.Get("/", func(ctx *tigo.Context) error {
        content := `
            Hello tigo!!!<hr>
            visit admin: <a href="/admin/">/admin/</a>
        `
        //out json
        return ctx.HTML(content)
    })

    //register admin router
    admin.Get("/", func(ctx *tigo.Context) error {
        return ctx.Render("admin/page", tigo.M{"title": "Tigo render"})
    })

    //run
    log.Printf("run on :8080")
    err := router.Run(":8080")
    if err != nil {
        log.Fatalf("run error: %v", err)
    }
	
```

/views/admin/page.html

```html

        <!-- /views/admin/page.html content -->

        {{layout "admin/master"}}

        <h3>admin/page.html</h3>
        <div>this admin/page.html</div>

```


/views/admin/master.html
```html
        <!-- /views/admin/master.html content -->

        <!doctype html>

        <html>
        <head>
            <meta http-equiv="Content-type" content="text/html; charset=utf-8" />
            <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
            <title>{{.title}}</title>
        </head>

        <body>
        admin/master.html

        <hr>
        render page content will at here:
        {{content}}
        </body>
        </html>
        
```


Now run the following command to start the Web server:

```
go run main.go
```

You should be able to access URLs such as `http://localhost:8080`.


### Routes

tigo works by building a routing table in a router and then dispatching HTTP requests to the matching handlers 
found in the routing table. An intuitive illustration of a routing table is as follows:


Routes              |  Handlers
--------------------|-----------------
`GET /users`        |  m1, m2, h1, ...
`POST /users`       |  m1, m2, h2, ...
`PUT /users/<id>`   |  m1, m2, h3, ...
`DELETE /users/<id>`|  m1, m2, h4, ...


For an incoming request `GET /users`, the first route would match and the handlers m1, m2, and h1 would be executed.
If the request is `PUT /users/123`, the third route would match and the corresponding handlers would be executed.
Note that the token `<id>` can match any number of non-slash characters and the matching part can be accessed as 
a path parameter value in the handlers.

**If an incoming request matches multiple routes in the table, the route added first to the table will take precedence.
All other matching routes will be ignored.**

The actual implementation of the routing table uses a variant of the radix tree data structure, which makes the routing
process as fast as working with a hash table, thanks to the inspiration from [httprouter](https://github.com/julienschmidt/httprouter).

To add a new route and its handlers to the routing table, call the `To` method like the following:
  
```go
router := tigo.New()
router.To("GET", "/users", m1, m2, h1)
router.To("POST", "/users", m1, m2, h2)
```

You can also use shortcut methods, such as `Get`, `Post`, `Put`, etc., which are named after the HTTP method names:
 
```go
router.Get("/users", m1, m2, h1)
router.Post("/users", m1, m2, h2)
```

If you have multiple routes with the same URL path but different HTTP methods, like the above example, you can 
chain them together as follows,

```go
router.Get("/users", m1, m2, h1).Post(m1, m2, h2)
```

If you want to use the same set of handlers to handle the same URL path but different HTTP methods, you can take
the following shortcut:

```go
router.To("GET,POST", "/users", m1, m2, h)
```

A route may contain parameter tokens which are in the format of `<name:pattern>`, where `name` stands for the parameter
name, and `pattern` is a regular expression which the parameter value should match. A token `<name>` is equivalent
to `<name:[^/]*>`, i.e., it matches any number of non-slash characters. At the end of a route, an asterisk character
can be used to match any number of arbitrary characters. Below are some examples:

* `/users/<username>`: matches `/users/admin`
* `/users/accnt-<id:\d+>`: matches `/users/accnt-123`, but not `/users/accnt-admin`
* `/users/<username>/*`: matches `/users/admin/profile/address`

When a URL path matches a route, the matching parameters on the URL path can be accessed via `Context.Param()`:

```go
router := tigo.New()

router.Get("/users/<username>", func (c *tigo.Context) error {
	fmt.Fprintf(c, "Name: %v", c.Param("username"))
	return nil
})
```


### Route Groups

Route group is a way of grouping together the routes which have the same route prefix. The routes in a group also
share the same handlers that are registered with the group via its `Use` method. For example,

```go
router := tigo.New()
api := router.Group("/api")
api.Use(m1, m2)
api.Get("/users", h1).Post(h2)
api.Put("/users/<id>", h3).Delete(h4)
```

The above `/api` route group establishes the following tigo table:


Routes                  |  Handlers
------------------------|-------------
`GET /api/users`        |  m1, m2, h1, ...
`POST /api/users`       |  m1, m2, h2, ...
`PUT /api/users/<id>`   |  m1, m2, h3, ...
`DELETE /api/users/<id>`|  m1, m2, h4, ...


As you can see, all these routes have the same route prefix `/api` and the handlers `m1` and `m2`. In other similar
tigo frameworks, the handlers registered with a route group are also called *middlewares*.

Route groups can be nested. That is, a route group can create a child group by calling the `Group()` method. The router
serves as the top level route group. A child group inherits the handlers registered with its parent group. For example, 

```go
router := tigo.New()
router.Use(m1)

api := router.Group("/api")
api.Use(m2)

users := group.Group("/users")
users.Use(m3)
users.Put("/<id>", h1)
```

Because the router serves as the parent of the `api` group which is the parent of the `users` group, 
the `PUT /api/users/<id>` route is associated with the handlers `m1`, `m2`, `m3`, and `h1`.


### Router

Router manages the tigo table and dispatches incoming requests to appropriate handlers. A router instance is created
by calling the `tigo.New()` method.

To hook up router with fasthttp, use the following code:

```go
router := tigo.New()
fasthttp.ListenAndServe(":8080", router.HandleRequest) 
```


### Handlers

A handler is a function with the signature `func(*tigo.Context) error`. A handler is executed by the router if
the incoming request URL path matches the route that the handler is associated with. Through the `tigo.Context` 
parameter, you can access the request information in handlers.

A route may be associated with multiple handlers. These handlers will be executed in the order that they are registered
to the route. The execution sequence can be terminated in the middle using one of the following two methods:

* A handler returns an error: the router will skip the rest of the handlers and handle the returned error.
* A handler calls `Context.Abort()`: the router will simply skip the rest of the handlers. There is no error to be handled.
 
A handler can call `Context.Next()` to explicitly execute the rest of the unexecuted handlers and take actions after
they finish execution. For example, a response compression handler may start the output buffer, call `Context.Next()`,
and then compress and send the output to response.


### Context

For each incoming request, a `tigo.Context` object is passed through the relevant handlers. Because `tigo.Context`
embeds `fasthttp.RequestCtx`, you can access all properties and methods provided by the latter.
 
Additionally, the `Context.Param()` method allows handlers to access the URL path parameters that match the current route.
Using `Context.Get()` and `Context.Set()`, handlers can share data between each other. For example, an authentication
handler can store the authenticated user identity by calling `Context.Set()`, and other handlers can retrieve back
the identity information by calling `Context.Get()`.

Context also provides a handy `WriteData()` method that can be used to write data of arbitrary type to the response.
The `WriteData()` method can also be overridden (by replacement) to achieve more versatile response data writing. 


### Error Handling

A handler may return an error indicating some erroneous condition. Sometimes, a handler or the code it calls may cause
a panic. Both should be handled properly to ensure best user experience. It is recommended that you use 
the `fault.Recover` handler or a similar error handler to handle these errors.

If an error is not handled by any handler, the router will handle it by calling its `handleError()` method which
simply sets an appropriate HTTP status code and writes the error message to the response.

When an incoming request has no matching route, the router will call the handlers registered via the `Router.NotFound()`
method. All the handlers registered via `Router.Use()` will also be called in advance. By default, the following two
handlers are registered with `Router.NotFound()`:

* `tigo.MethodNotAllowedHandler`: a handler that sends an `Allow` HTTP header indicating the allowed HTTP methods for a requested URL
* `tigo.NotFoundHandler`: a handler triggering 404 HTTP error


## Docs

See http://godoc.org/github.com/foolin/tigo .

