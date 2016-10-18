package main

import (
	"github.com/foolin/tigo"
	"log"
)

func main()  {

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
}