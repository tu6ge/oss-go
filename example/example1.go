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

	err = obj.Upload(content, "text/plain;charset=utf-8", &client)
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
	err = obj_copy.CopyFrom("/honglei123/aaabbc.html", "text/plain;charset=utf-8", &client)
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
