package main

import (
	"mygee2/gee"
)

func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.String(200, "URL.Path = %q\n", c.Path)
	})

	r.GET("/hello", func(c *gee.Context) {
		c.JSON(200, map[string]any{"name": "Tom", "age": 12})
	})
	r.Run(":9999")
}
