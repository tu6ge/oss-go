package main

import (
	"fmt"

	"github.com/tu6ge/oss-go"
)

// 分片下载的示例
func main() {
	// 初始化 client
	client, err := oss.NewWithEnv()
	if err != nil {
		fmt.Println(err)
		return
	}

	object := oss.NewPartsDownload("video222.mov")

	err = object.FilePath("./video.mov").Download(&client)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("download success")
}
