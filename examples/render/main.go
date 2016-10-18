package main

import (
	"github.com/foolin/tigo"
	"log"
)

func main()  {

	//new router
	router := tigo.New()

	//set render, tigo.Default() will default initialize.
	router.SetRender(tigo.NewHtmlRender(tigo.HtmlRenderConfig{
		ViewRoot:  "views",
		Extension: ".html",
	}))

	//register router
	router.Get("/", func(ctx *tigo.Context) error {
		/*
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
		*/

		/*
		    <!-- /views/layout/footer.html content -->

		    Copyright &copy2016 by <a href="https://github.com/foolin/tigo">tigo</a>.

		*/
		return ctx.Render("page", tigo.M{"title": "Tigo render"})
	})

	//run
	log.Printf("run on :8080")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("run error: %v", err)
	}
}