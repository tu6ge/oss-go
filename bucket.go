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
	query    types.ObjectQuery
	domain   string
}

func NewBucket(name, endpoint string) (Bucket, error) {
	end, err := types.NewEndPoint(endpoint)
	if err != nil {
		return Bucket{}, err
	}

	if len(name) == 0 {
		return Bucket{}, &InvalidBucketName{}
	}

	return Bucket{name, end, types.NewObjectQuery(), ""}, nil
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

	return Bucket{name, end, types.NewObjectQuery(), ""}, err
}

func (b *Bucket) SetEndPointDomain(domain string) error {
	return b.endpoint.SetOriginalDomain(domain)
}

func (b *Bucket) SetDomain(domain string) {
	b.domain = domain
}

func (b *Bucket) ToUrl() url.URL {
	u, _ := url.Parse("https://" + b.name + "." + b.endpoint.Host())

	if len(b.domain) > 0 {
		u, _ = url.Parse(b.domain)
	}

	return *u
}

func (b Bucket) Query(query map[string]string) Bucket {
	for key, val := range query {
		b.query.Insert(key, val)
	}
	return b
}

func (b Bucket) ObjectQuery(query types.ObjectQuery) Bucket {
	b.query = query
	return b
}

func (b Bucket) GetObjects(client *Client) (Objects, error) {
	url := b.ToUrl()
	url.RawQuery = b.query.ToOssQuery()
	method := "GET"
	resource := NewCanonicalizedResourceFromObjects(&b, b.query.GetNextToken())

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

	body_string := string(body)

	// fmt.Println("status", resp.StatusCode)

	if http_status_ok(resp.StatusCode) {
		token := parse_item(body_string, "NextContinuationToken")
		object_rs := parser_xml_objects(body_string)

		return Objects{object_rs, token, b.query}, nil
	} else {
		// fmt.Println(body_string)
		return Objects{}, parse_oss_response_error(body_string)
	}
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

		buckets = append(buckets, NewObject(name))
	}
	return buckets
}

func parse_item(xml, field string) string {
	start_tag := fmt.Sprintf("<%s>", field)
	end_tag := fmt.Sprintf("</%s>", field)
	start_index := strings.Index(xml, start_tag)
	end_index := strings.Index(xml, end_tag)

	if start_index > -1 && end_index > -1 {
		return xml[start_index+len(start_tag) : end_index]
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
