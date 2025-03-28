# Go 实现的 aliyun oss sdk

## 用法

1. 在项目根目录创建 `.env` 文件，在文件内填写 oss 的配置信息：
```
ALIYUN_KEY_ID=xxx
ALIYUN_KEY_SECRET=xxx
ALIYUN_ENDPOINT=cn-shanghai
ALIYUN_BUCKET=xxx
```

2. 运行如下代码
```go
package main

import (
	"fmt"

	"github.com/tu6ge/oss-go"
	"github.com/tu6ge/oss-go/types"
)

func main() {
	// Get a greeting message and print it.
	client, err := oss.NewWithEnv()
	if err != nil {
		fmt.Println(err)
		return
	}

	end, err := types.NewEndPoint("cn-shanghai")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(end)

	buckets, err := client.GetBuckets(end)
	if err != nil {
		fmt.Println(err)
		return
	}

	query := types.NewObjectQuery()
	query.Insert(types.QUERY_MAX_KEYS, "5")

	objects, _ := buckets[1].GetObjects(query, &client)

	fmt.Println(objects)

	obj := oss.NewObject("aaabbc4.html")

	content := []byte("foo")

	err = obj.Upload(content, "text/plain;charset=utf-8", &client)
	fmt.Println(err)

	con, err := obj.Download(&client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("content:", string(con))

	obj_copy := oss.NewObject("xyz.html")
	res := obj_copy.CopyFrom("/honglei123/aaabbc.html", "text/plain;charset=utf-8", &client)
	fmt.Println(res)

	err = obj.Delete(&client)

	fmt.Println(err)
}
```