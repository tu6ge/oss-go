package main

import (
	"fmt"

	"github.com/tu6ge/oss-go"
)

func main() {
	// 初始化 client
	client, err := oss.NewWithEnv()
	if err != nil {
		fmt.Println(err)
		return
	}

	object := oss.NewPartsUpload("video222.mov")

	err = object.FilePath("./video.mov").Upload(&client)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("upload success")
}
