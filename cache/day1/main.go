package main

import (
	"fmt"
	"gee_cache1/lru"
)

type String string

func (s String) Len() int {
	return len(s)
}
func main() {
	cache := lru.New(30)
	cache.Add("hello1", String("hello1"))
	cache.Add("hello2", String("hello2"))
	cache.Get("hello1")
	cache.Add("hello3", String("hello3"))
	fmt.Println("hello1", cache.Get("hello1"))
	fmt.Println("hello2", cache.Get("hello2"))
	cache.Add("hello4", String("hello4"))
	cache.Get("hello1")
	cache.Add("hello5", String("hello5"))

	fmt.Println("hello1", cache.Get("hello1"))
}
