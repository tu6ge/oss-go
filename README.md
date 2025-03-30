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
	"os"

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

	objects, err := buckets[1].GetObjects(query, &client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(objects)

	second_objects, err := objects.NextList(query, &client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(second_objects)

	obj := oss.NewObject("aaabbc4.html")

	content := []byte("foo")

	err = obj.Content(content).ContentType("text/plain;charset=utf-8").Upload(&client)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Open("./demofile.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	err = oss.NewObject("from_file.txt").File(f).ContentType("text/plain;charset=utf-8").Upload(&client)
	if err != nil {
		fmt.Println(err)
		return
	}

	con, err := obj.Download(&client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("content:", string(con))

	obj_copy := oss.NewObject("xyz.html")
	err = obj_copy.CopySource("/honglei123/aaabbc.html").ContentType("text/plain;charset=utf-8").Copy(&client)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = obj.Delete(&client)
	if err != nil {
		fmt.Println(err)
		return
	}
}
```

# Bench

跟 aliyun 官方提供的 sdk 进行 bench 比较，发现性能提高了不少，以下是上传文件进行 bench 的测试

```
本 library           30236             39481 ns/op
aliyun官方的sdk          10         100343922 ns/op
```
