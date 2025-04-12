package oss

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/tu6ge/oss-go/types"

	"github.com/joho/godotenv"
)

var (
	QUERY_START_AFTER        = types.QUERY_START_AFTER
	QUERY_CONTINUATION_TOKEN = types.QUERY_CONTINUATION_TOKEN
	QUERY_MAX_KEYS           = types.QUERY_MAX_KEYS
	QUERY_PREFIX             = types.QUERY_PREFIX
	QUERY_ENCODING_TYPE      = types.QUERY_ENCODING_TYPE
	QUERY_FETCH_OWNER        = types.QUERY_FETCH_OWNER
)

type Client struct {
	access_key_id    string
	access_secret_id types.Secret
	Bucket           Bucket
}

func New(key, secret, bucket, endpoint string) (Client, error) {
	bucket_name, err := NewBucket(bucket, endpoint)
	if err != nil {
		return Client{}, err
	}
	return Client{
		key,
		types.NewSecret(secret),
		bucket_name,
	}, nil
}

func NewWithEnv() (Client, error) {
	err := godotenv.Load()
	if err != nil {
		return Client{}, err
	}

	// 读取环境变量
	key_id := os.Getenv("ALIYUN_KEY_ID")
	secret_id := os.Getenv("ALIYUN_KEY_SECRET")

	if key_id == "" || secret_id == "" {
		return Client{}, &EnvEmtpyError{}
	}

	bucket, err := BucketFromEnv()
	if err != nil {
		return Client{}, err
	}

	return Client{key_id, types.NewSecret(secret_id), bucket}, nil
}

func (c Client) Authorization(method string, resource types.CanonicalizedResource) map[string]string {
	return c.AuthorizationHeader(method, resource, make(map[string]string, 0))
}

const (
	LINE_BREAK   string = "\n"
	CONTENT_TYPE string = "text/xml"
)

func (c Client) AuthorizationHeader(method string, resource types.CanonicalizedResource, headers map[string]string) map[string]string {
	date := now()

	resource_str := resource.ToStr()

	oss_header_str := to_oss_header(headers)

	content_type, ok_content_type := headers["Content-Type"]

	sign_str := method
	sign_str += LINE_BREAK
	sign_str += LINE_BREAK
	if ok_content_type {
		sign_str += content_type
	}
	sign_str += LINE_BREAK
	sign_str += date
	sign_str += LINE_BREAK
	sign_str += oss_header_str
	sign_str += resource_str

	// fmt.Println("sign_str", sign_str)

	encry := c.access_secret_id.Encryption(sign_str)

	sign := fmt.Sprintf("OSS %s:%s", c.access_key_id, encry)

	headers["AccessKeyId"] = c.access_key_id
	headers["VERB"] = method
	headers["Date"] = date
	headers["Authorization"] = sign
	headers["CanonicalizedResource"] = resource_str
	return headers
}

// func (c *Client) SetEndPointDomain(domain string) error {
// 	return c.Bucket.endpoint.SetOriginalDomain(domain)
// }

func (c *Client) SetBucketDomain(domain string) {
	c.Bucket.SetDomain(domain)
}

func (c Client) GetBuckets(endpoint ...types.EndPoint) ([]Bucket, error) {
	var url url.URL
	var end types.EndPoint
	if len(endpoint) == 0 {
		end = c.Bucket.endpoint
		url = end.ToUrl()
	} else if len(endpoint) == 1 {
		end = endpoint[0]
		url = end.ToUrl()
	} else {
		return []Bucket{}, errors.New("too many args")
	}
	method := "GET"
	resource := types.DefaultCanonicalizedResource()

	header_map := c.Authorization(method, resource)

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return []Bucket{}, err
	}

	for k, v := range header_map {
		req.Header.Add(k, v)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return []Bucket{}, err
	}

	defer resp.Body.Close() // 关闭响应体

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Bucket{}, err
	}

	body_string := string(body)

	if http_status_ok(resp.StatusCode) {
		return parser_xml(body_string, end), nil
	} else {
		return nil, parse_oss_response_error(body_string)
	}
}

func http_status_ok(status int) bool {
	return status >= http.StatusOK && status < http.StatusMultipleChoices
}

func parser_xml(xml string, endpoint types.EndPoint) []Bucket {
	var start_positions []int
	var end_positions []int
	start := 0
	pattern := "<Name>"
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
	pattern = "</Name>"
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

	var buckets []Bucket
	for i, item := range start_positions {
		name := xml[item+len("<Name>") : end_positions[i]]

		buckets = append(buckets, Bucket{name, endpoint, types.NewObjectQuery(), ""})
	}
	return buckets
}

func now() string {
	// 获取当前时间并转换为 UTC
	currentTime := time.Now().UTC()

	// 格式化时间为 RFC1123 格式（带 GMT）
	formattedTime := currentTime.Format(time.RFC1123)

	// 将 "UTC" 替换为 "GMT"
	formattedTime = replaceUTCWithGMT(formattedTime)
	return formattedTime
}

// 替换 UTC 为 GMT
func replaceUTCWithGMT(timeStr string) string {
	if len(timeStr) > 3 && timeStr[len(timeStr)-3:] == "UTC" {
		return timeStr[:len(timeStr)-3] + "GMT"
	}
	return timeStr
}

type ossMap struct {
	key   string
	value string
}

func to_oss_header(headers map[string]string) string {
	var oss_map []ossMap
	for k, v := range headers {
		if strings.HasPrefix(k, "x-oss-") {
			oss_map = append(oss_map, ossMap{k, v})
		}
	}

	if len(oss_map) == 0 {
		return ""
	}

	sort.Slice(oss_map, func(i, j int) bool {
		return oss_map[i].key < oss_map[j].key
	})

	result := ""
	for _, item := range oss_map {
		result += item.key + ":" + item.value + "\n"
	}

	return result
}
