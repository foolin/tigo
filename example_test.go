package tigo_test

import (
	"github.com/foolin/tigo"
	"log"
	"os"
)

//Example default example
func Example() {
	//new default router
	router := tigo.Default()

	//register router
	router.Get("/api/<action>", func(ctx *tigo.Context) error {

		//json object
		data := struct {
			Ip string
			Action string
		}{ctx.RequestIP(), ctx.Param("action")}

		//out json
		return ctx.JSON(data)
	})

	//run
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("run error: %v", err)
	}
}

//Example_New example new
func Example_new() {

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
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("run error: %v", err)
	}
}



//Example_New render html
//
func Example_render() {

	//new router
	router := tigo.New()

	//tigo.Default() init like this.
	router.SetRender(tigo.NewHtmlRender(tigo.HtmlRenderConfig{
		ViewRoot: "views",
		Extension: ".html",
	}))

	//register router
	router.Get("/", func(ctx *tigo.Context) error {
		/*
		<---- /views/page.html content --->

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
		<---- /views/layout/footer.html content --->

		@copryright 2016 by <a href="https://github.com/foolin/tigo">tigo</a>.

		 */
		return ctx.Render("page", tigo.M{"title", "Tigo render"})
	})

	//run
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("run error: %v", err)
	}
}


//Example_RenderMaster render html user master template.
//
func Example_renderMaster() {

	//new router
	router := tigo.Default()

	admin := router.Group("/admin")

	//register router
	admin.Get("/", func(ctx *tigo.Context) error {

		/*

		<---- /views/admin/page.html content --->

		{{layout "admin/master"}}

		<div>this admin/page.html</div>

		*/

		/*

		<---- /views/admin/master.html content --->

		<!doctype html>

		<html>
		<head>
		    <meta http-equiv="Content-type" content="text/html; charset=utf-8" />
		    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
		    <title>{{.title}}</title>
		</head>

		<body>
			admin/page.html content will at here:
			{content}
		</body>
		</html>

		 */
		return ctx.Render("admin/page", tigo.M{"title", "Tigo render"})
	})

	//run
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("run error: %v", err)
	}
}