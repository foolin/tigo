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
		return ctx.Render("admin/page", tigo.M{"title": "Tigo render"})
	})

	//run
	log.Printf("run on :8080")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("run error: %v", err)
	}
}
