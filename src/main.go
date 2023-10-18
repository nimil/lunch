package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {

	lunchArr := []string{"黄焖鸡", "食堂炒菜", "面条", "超意兴"}
	rand.NewSource(time.Now().UnixNano())
	randomLunch := lunchArr[rand.Intn(len(lunchArr))]

	fmt.Fprintf(w, "<h1><center>"+randomLunch+"</center></h1>") //这个写入到w的是输出到客户端的
}

func main() {
	http.HandleFunc("/", sayhelloName)       //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
