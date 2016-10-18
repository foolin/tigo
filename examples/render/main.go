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
		return ctx.Render("page", tigo.M{"title": "Tigo render"})
	})

	//run
	log.Printf("run on :8080")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("run error: %v", err)
	}
}