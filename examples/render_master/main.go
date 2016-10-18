package main

import (
	"github.com/foolin/tigo"
	"log"
)

func main()  {
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

		/*

		<!-- /views/admin/page.html content -->

		{{layout "admin/master"}}

		<h3>admin/page.html</h3>
		<div>this admin/page.html</div>

		*/

		/*

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

		 */
		return ctx.Render("admin/page", tigo.M{"title": "Tigo render"})
	})

	//run
	log.Printf("run on :8080")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("run error: %v", err)
	}
}
