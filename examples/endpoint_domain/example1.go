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
	// 或者
	// client,err := oss.New("key","secret","bucket_name","cn-hangzhou")
	if err != nil {
		fmt.Println(err)
		return
	}

	// client.SetEndPointDomain("https://oss-accelerate.aliyuncs.com")

	// 根据默认配置的 endpoint 获取 bucket 列表
	buckets_from_cofig, err := client.GetBuckets()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(buckets_from_cofig)

	end, err := types.NewEndPoint("cn-shanghai")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(end)

	// 获取所有 bucket
	buckets, err := client.GetBuckets(end)
	if err != nil {
		fmt.Println(err)
		return
	}

	buckets[1].SetEndPointDomain("https://oss-accelerate.aliyuncs.com")

	// 查询文件列表
	query := map[string]string{
		oss.QUERY_MAX_KEYS: "5",
	}

	objects, err := buckets[1].Query(query).GetObjects(&client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("objects :", objects)

	// 查询第二页的文件列表
	second_objects, err := objects.NextList(&client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("second objects :", second_objects)

	// 初始化文件结构体
	obj := oss.NewObject("aaabbc4.html")

	client.SetBucketDomain("http://oss2.xiaoqiezi.top")

	// 使用文件内容上传文件
	content := []byte("foo")

	err = obj.Content(content).ContentType("text/plain;charset=utf-8").Upload(&client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 使用文件句柄上传文件
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

	// 使用本地文件路径上传文件
	err = oss.NewObject("from_file2.txt").FilePath("./demofile.txt").ContentType("text/plain;charset=utf-8").Upload(&client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 下载文件内容
	con, err := obj.Download(&client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("content:", string(con))

	// 复制文件
	obj_copy := oss.NewObject("xyz.html")
	err = obj_copy.CopySource("/honglei123/aaabbc.html").ContentType("text/plain;charset=utf-8").Copy(&client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 删除文件
	err = obj.Delete(&client)
	if err != nil {
		fmt.Println(err)
		return
	}
}
