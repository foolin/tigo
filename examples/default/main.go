package main

import (
	"github.com/foolin/tigo"
	"log"
)

func main()  {

	router := tigo.Default()

	router.Get("/", func(ctx *tigo.Context) error {
		//out html
		return ctx.HTML(`
			Hello tigo!!!<hr>
			visit api: <a href="/api/tigo">api/tigo</a>
		`)
	})

	router.Get("/api/<action>", func(ctx *tigo.Context) error {
		//out json
		return ctx.JSON(tigo.M{
			"name": "tigo",
			"ip": ctx.RequestIP(),
			"action": ctx.Param("action"),
		})
	})

	//run
	log.Printf("run on :8080")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("run error: %v", err)
	}
}