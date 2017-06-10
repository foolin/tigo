package main

import (
	"github.com/foolin/tigo"
	"fmt"
)

func main()  {
	router := tigo.Default()

	router.NotFound(func(ctx *tigo.Context) error {
		//out html
		return ctx.HTML(`
			^_^!!! 404 Not Found!
		`)
	})

	router.OnError = func(ctx *tigo.Context, err error) {
		//out html
		ctx.HTML(fmt.Sprintf("^_^!!! 500 Server Error! %v", err.Error()))
	}

	router.Get("/err", func(ctx *tigo.Context) error {
		a := 1
		b := 0
		c := a / b
		print(c)
		return nil
	})

	//run
	println("run on :8080")
	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
