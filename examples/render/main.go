package main

import (
	"github.com/foolin/tigo"
	"log"
	"html/template"
)

func main() {

	//new router
	router := tigo.New()

	//set render, tigo.Default() will default initialize.
	router.SetRender(tigo.NewViewRender(tigo.ViewRenderConfig{
		Root: "views",
		Extension: ".html",
		Master: "layout/master",
		Partials: []string{"layout/footer"},
		Funcs: template.FuncMap{
			"echo": func(content string) template.HTML {
				return template.HTML(content)
			},
		},
		DisableCache: false,
		EnableFilePartial: true,
	}))

	//register router
	router.Get("/", func(ctx *tigo.Context) error {
		return ctx.Render("index", tigo.M{
			"title": "Index title!",
			"escape": func(content string) string {
				return template.HTMLEscapeString(content)
			},
		})
	})

	router.Get("/page_file", func(ctx *tigo.Context) error {
		return ctx.RenderFile("page_file", tigo.M{"title": "Page file title!!"})
	})

	//run
	log.Printf("run on :8080")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("run error: %v", err)
	}
}