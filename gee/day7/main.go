package main

import (
	"fmt"
	"mygee7/gee"
)

func main() {
	router := gee.New()
	router.Use(gee.Recovery())
	router.GET("/panic", func(context *gee.Context) {
		str := []string{"hello", "world"}
		//数组越界
		fmt.Println(str[12])
	})
	router.Run((":9999"))
}
