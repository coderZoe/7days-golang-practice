package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	//说明：这里我们配置了路径，但奇怪的是似乎没将这些路由传给ListenAndServe
	//这是因为通过http.HandleFunc配置路径时，golang会默认创建一个全局的http.DefaultServeMux来装你配置的路由
	//当ListenAndServe第二个参数传nil时，会使用全局的http.DefaultServeMux作为路由表
	//我们当然也可以自己创建ServeMux
	//自己创建的版本见 main2函数
	//一般建议自己实现，因为全局的DefaultServeMux会可能被篡改
	//另外还需要补充的是：ListenAndServe接收的第二个参数是Handler
	//ServeMux是一个Handler，其中它内部包含了多个我们注册给它的Handler
	//ServeMux.Handle()/http.Handle()方法接收的是一个Handler
	//而ServeMux.HandleFunc和http.HandleFunc接收的是一个function，但在内部会通过HandlerFunc将function转为Handler
	//详见 https://www.alexedwards.net/blog/an-introduction-to-handlers-and-servemuxes-in-go
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func indexHandler(writer http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf函数，第一个参数是写出流，第二个参数是字符串
	fmt.Fprintf(writer, "URL.Path =%q\n", req.URL.Path)
}

func helloHandler(writer http.ResponseWriter, req *http.Request) {
	for key, value := range req.Header {
		fmt.Fprintf(writer, "Header[%q]=%q\n", key, value)
	}
}

func main2() {
	serveMus := http.NewServeMux()
	serveMus.HandleFunc("/", indexHandler)
	serveMus.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":9999", serveMus))
}
