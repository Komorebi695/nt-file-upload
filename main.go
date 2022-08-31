package main

import (
	"fmt"
	"net/http"
	"ntfileupload/handler"
)

func main() {
	// 静态资源处理
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./view"))))

	// 文件上传
	http.HandleFunc("/nt/upload", handler.UploadHandler)

	// 监听端口
	err := http.ListenAndServe(":8082", nil)
	if nil != err {
		fmt.Printf("Failed to start server,err:%s\n", err.Error())
	}

}
