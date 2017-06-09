package main

import (
	"gopkg.in/foolin/tigo.v1"
	"os"
	"log"
)

func main()  {
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
}
