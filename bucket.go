package oss

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/tu6ge/oss-go/types"
)

type Bucket struct {
	name     string
	endpoint types.EndPoint
}

func NewBucket(name, endpoint string) (Bucket, error) {
	end, err := types.NewEndPoint(endpoint)
	if err != nil {
		return Bucket{}, err
	}

	if len(name) == 0 {
		return Bucket{}, &InvalidBucketName{}
	}

	return Bucket{name, end}, nil
}

func BucketFromEnv() (Bucket, error) {
	err := godotenv.Load()
	if err != nil {
		return Bucket{}, err
	}
	name := os.Getenv("ALIYUN_BUCKET")
	end, err := types.EndPointFromEnv()
	if err != nil {
		return Bucket{}, err
	}

	return Bucket{name, end}, err
}

func (b *Bucket) ToUrl() url.URL {
	u, _ := url.Parse("https://" + b.name + "." + b.endpoint.Host())

	return *u
}

func (b *Bucket) GetObjects(query types.ObjectQuery, client Client) (Objects, error) {
	url := b.ToUrl()
	url.RawQuery = query.ToOssQuery()
	method := "GET"
	resource := NewCanonicalizedResourceFromObjects(b, query.GetNextToken())

	headers := client.Authorization(method, resource)

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return Objects{}, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	http_client := http.DefaultClient
	resp, err := http_client.Do(req)
	if err != nil {
		return Objects{}, err
	}

	defer resp.Body.Close() // 关闭响应体

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Objects{}, err
	}

	// fmt.Println(string(body))
	xml := string(body)
	token := parse_item(xml, "NextContinuationToken")
	object_rs := parser_xml_objects(xml)

	return Objects{object_rs, token}, nil
}

func parser_xml_objects(xml string) []Object {
	var start_positions []int
	var end_positions []int
	start := 0
	pattern := "<Key>"
	pattern_len := len(pattern)

	pos := strings.Index(xml[start:], pattern)
	for pos > -1 {
		start_positions = append(start_positions, pos)
		start = pos + pattern_len
		if start > len(xml) {
			break
		}
		pos = strings.Index(xml[start:], pattern)
		if pos == -1 {
			break
		}
		pos = pos + start
	}

	start = 0
	pattern = "</Key>"
	pattern_len = len(pattern)
	pos = strings.Index(xml[start:], pattern)
	for pos > -1 {
		end_positions = append(end_positions, pos)
		start = pos + pattern_len
		if start > len(xml) {
			break
		}
		pos = strings.Index(xml[start:], pattern)
		if pos == -1 {
			break
		}
		pos = pos + start
	}

	var buckets []Object
	for i, item := range start_positions {
		name := xml[item+len("<Key>") : end_positions[i]]

		buckets = append(buckets, Object{name})
	}
	return buckets
}

func parse_item(xml, field string) string {
	start_tag := fmt.Sprintf("<%s>", field)
	end_tag := fmt.Sprintf("</%s>", field)
	start_index := strings.Index(xml, start_tag)
	end_index := strings.Index(xml, end_tag)

	if start_index > -1 && end_index > -1 {
		return xml[start_index+2 : end_index]
	}
	return ""
}

func NewCanonicalizedResourceFromObjects(bucket *Bucket, continuation_token string) types.CanonicalizedResource {
	if len(continuation_token) > 0 {
		return types.NewCanonicalizedResource(fmt.Sprintf("/%s/?continuation-token=%s", bucket.name, continuation_token))
	}
	return types.NewCanonicalizedResource(fmt.Sprintf("/%s/", bucket.name))
}

type InvalidBucketName struct{}

func (b InvalidBucketName) Error() string {
	return "invalid bucket name"
}
