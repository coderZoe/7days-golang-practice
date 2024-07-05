package main

import "mygee4/gee"

func main() {
	//我们的e既是engine 也是routerGroup
	e := gee.New()
	adminGroup := e.Group("/admin")
	adminGroup.AddRouter("GET", "/info", func(ctx *gee.Context) {
		ctx.JSON(200, map[string]any{
			"name": "admin",
			"age":  12,
		})
	})

	adminGroup.AddRouter("POST", "/update/:name", func(ctx *gee.Context) {
		ctx.JSON(200, map[string]any{
			"name": ctx.Params["name"],
			"age":  12,
		})
	})

	adminAccountRg := adminGroup.Group("/account")
	adminAccountRg.GET("/info", func(ctx *gee.Context) {
		ctx.JSON(200, map[string]any{
			"salary": 2800,
			"left":   23.5,
		})
	})
	e.Run(":9999")
}
