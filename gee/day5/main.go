package main

import (
	"log"
	"mygee5/gee"
	"time"
)

func main() {
	//我们的e既是engine 也是routerGroup
	e := gee.New()
	v1 := e.Group("/v1")
	v1.Use(func(ctx *gee.Context) {
		path := ctx.Path
		method := ctx.Method
		log.Printf("v1 middleware method :%s,path:%s,start at %q\n", method, path, time.Now())
		ctx.Next()
		log.Printf("v1 middleware method :%s,path:%s,end at %q\n", method, path, time.Now())
	})

	v1.AddRouter("GET", "/info", func(ctx *gee.Context) {
		ctx.JSON(200, map[string]any{
			"name": "admin",
			"age":  12,
		})
	})

	v1Admin := v1.Group("/admin")
	v1Admin.Use(func(ctx *gee.Context) {
		path := ctx.Path
		method := ctx.Method
		log.Printf("v1/admin middleware method :%s,path:%s,start at %q\n", method, path, time.Now())
		ctx.Next()
		log.Printf("v1/admin middleware method :%s,path:%s,end at %q\n", method, path, time.Now())
	})

	v1Admin.AddRouter("GET", "/update/:name", func(ctx *gee.Context) {
		ctx.JSON(200, map[string]any{
			"name": ctx.Params["name"],
			"age":  12,
		})
	})

	e.Run(":9999")
}
