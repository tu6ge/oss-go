package types

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"os"
	"strings"
	"unicode"

	"github.com/joho/godotenv"
)

type Secret struct {
	value string
}

func NewSecret(value string) Secret {
	return Secret{value}
}

func (s Secret) Encryption(data string) string {
	// 将 key 转成字节数组
	secretKey := []byte(s.value)
	// 创建 HMAC-SHA1 对象
	h := hmac.New(sha1.New, secretKey)
	// 将数据写入哈希
	h.Write([]byte(data))
	// 计算最终的哈希值
	hashedData := h.Sum(nil)
	// 将哈希值转成十六进制字符串
	return base64.StdEncoding.EncodeToString(hashedData)
}

type CanonicalizedResource struct {
	value string
}

func NewCanonicalizedResource(value string) CanonicalizedResource {
	return CanonicalizedResource{value}
}

func DefaultCanonicalizedResource() CanonicalizedResource {
	return CanonicalizedResource{"/"}
}

func (c CanonicalizedResource) ToStr() string {
	return c.value
}

type ObjectQuery struct {
	query map[string]string
}

const (
	QUERY_DELIMITER          string = "delimiter"
	QUERY_START_AFTER        string = "start-after"
	QUERY_CONTINUATION_TOKEN string = "continuation-token"
	QUERY_MAX_KEYS           string = "max-keys"
	QUERY_PREFIX             string = "prefix"
	QUERY_ENCODING_TYPE      string = "encoding-type"
	QUERY_FETCH_OWNER        string = "fetch-owner"
)

func NewObjectQuery() ObjectQuery {
	return ObjectQuery{query: make(map[string]string)}
}

func (q ObjectQuery) Insert(key, value string) {
	q.query[key] = value
}

// TODO 处理 error
func (q ObjectQuery) GetNextToken() string {
	v, ok := q.query[QUERY_CONTINUATION_TOKEN]
	if ok {
		return v
	}
	return ""
}

func (q ObjectQuery) ToOssQuery() string {
	query_str := "list-type=2"
	for key, value := range q.query {
		query_str += "&"
		query_str += key
		query_str += "="
		query_str += value
	}
	return query_str
}

func (q ObjectQuery) Insert_next_token(value string) {
	q.query[QUERY_CONTINUATION_TOKEN] = value
}

type EndPoint struct {
	value       string
	is_internal bool
}

const (
	ENDPOINT_BEIJING     string = "cn-beijing"
	ENDPOINT_SHANGHAI    string = "cn-shanghai"
	ENDPOINT_QINGDAO     string = "cn-qingdao"
	ENDPOINT_SHENZHEN    string = "cn-shenzhen"
	ENDPOINT_HANGZHOU    string = "cn-hangzhou"
	ENDPOINT_HONGKONG    string = "cn-hongkong"
	ENDPOINT_GUANGZHOU   string = "cn-guangzhou"
	ENDPOINT_CHENGDU     string = "cn-chengdu"
	ENDPOINT_ZHANGJIAKOU string = "cn-zhangjiakou"
	ENDPOINT_HEFEI       string = "cn-hefei"
	ENDPOINT_WUHAN       string = "cn-wuhan"
	ENDPOINT_NANJING     string = "cn-nanjing"
	ENDPOINT_US_WEST_1   string = "us-west-1"
	ENDPOINT_US_EAST_1   string = "us-east-1"
)

func EndPointFromEnv() (EndPoint, error) {
	err := godotenv.Load()
	if err != nil {
		return EndPoint{}, err
	}
	str := os.Getenv("ALIYUN_ENDPOINT")

	return NewEndPoint(str)
}

func NewEndPoint(value string) (EndPoint, error) {
	if len(value) == 0 {
		return EndPoint{}, &InvalidEndPoint{}
	}

	is_internal := false

	if strings.HasSuffix(value, "-internal") {
		is_internal = true
		value = strings.Replace(value, "-internal", "", 1)
	}

	if strings.HasPrefix(value, "-") || strings.HasSuffix(value, "-") || strings.HasPrefix(value, "oss") {
		return EndPoint{}, &InvalidEndPoint{}
	}

	if !isValidString(value) {
		return EndPoint{}, &InvalidEndPoint{}
	}

	return EndPoint{value, is_internal}, nil
}

func (e EndPoint) SetInternal(is_internal bool) {
	e.is_internal = is_internal
}

func (e EndPoint) IsInternal() bool {
	return e.is_internal
}

func (e *EndPoint) ToUrl() url.URL {
	u, _ := url.Parse("https://" + e.Host())

	return *u
}

func (e *EndPoint) Host() string {
	host := "oss-"
	host += e.value
	if e.is_internal {
		host += "-internal"
	}
	host += ".aliyuncs.com"

	return host
}

func isValidString(s string) bool {
	for _, char := range s {
		if !(unicode.IsLetter(char) || unicode.IsDigit(char) || char == '-') {
			return false
		}
	}
	return true
}

func DefaultEndPoint() EndPoint {
	return EndPoint{ENDPOINT_QINGDAO, false}
}

type InvalidEndPoint struct{}

func (e *InvalidEndPoint) Error() string {
	return "invalid endpoint"
}
