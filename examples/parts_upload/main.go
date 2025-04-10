package main

import (
	"fmt"
	"os"

	"github.com/tu6ge/oss-go"
	"github.com/tu6ge/oss-go/types"
)

func main() {
	// 初始化 client
	client, err := oss.NewWithEnv()
	if err != nil {
		fmt.Println(err)
		return
	}
}
