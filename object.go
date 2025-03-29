package oss

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tu6ge/oss-go/types"
)

type Objects struct {
	List      []Object
	NextToken string
}

type Object struct {
	path string
}

func NewObject(path string) Object {
	return Object{path}
}

func (obj Object) ToUrl(bucket *Bucket) url.URL {
	url := bucket.ToUrl()
	url.Path = obj.path
	return url
}

func (obj Object) Upload(content []byte, content_type string, client *Client) error {
	bucket := client.bucket
	url := obj.ToUrl(&bucket)
	method := "PUT"

	resource := CanonicalizedResourceFromObject(&bucket, &obj)
	headers := make(map[string]string)
	if len(content_type) > 0 {
		headers["Content-Type"] = content_type
	}
	headers = client.AuthorizationHeader(method, resource, headers)

	if len(content) == 0 {
		headers["Content-Length"] = "0"
	}

	req, err := http.NewRequest(method, url.String(), bytes.NewReader(content))
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if http_status_ok(resp.StatusCode) {
		return nil
	} else {
		// 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		body_string := string(body)
		return parse_oss_response_error(body_string)
	}
}

func (obj Object) Download(client *Client) ([]byte, error) {
	bucket := client.bucket
	url := obj.ToUrl(&bucket)
	method := "GET"

	resource := CanonicalizedResourceFromObject(&bucket, &obj)

	headers := client.Authorization(method, resource)

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if http_status_ok(resp.StatusCode) {
		return data, nil
	} else {
		body_string := string(data)
		return nil, parse_oss_response_error(body_string)
	}
}

func (obj Object) CopyFrom(source string, content_type string, client *Client) error {
	bucket := client.bucket
	url := obj.ToUrl(&bucket)
	method := "PUT"

	resource := CanonicalizedResourceFromObject(&bucket, &obj)
	headers := make(map[string]string)
	headers["x-oss-copy-source"] = source
	if len(content_type) > 0 {
		headers["Content-Type"] = content_type
	}
	headers = client.AuthorizationHeader(method, resource, headers)

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if http_status_ok(resp.StatusCode) {
		return nil
	} else {
		// 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		body_string := string(body)
		return parse_oss_response_error(body_string)
	}
}

func (obj Object) Delete(client *Client) error {
	bucket := client.bucket
	url := obj.ToUrl(&bucket)
	method := "DELETE"

	resource := CanonicalizedResourceFromObject(&bucket, &obj)
	headers := client.Authorization(method, resource)

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if http_status_ok(resp.StatusCode) {
		return nil
	} else {
		// 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		body_string := string(body)
		return parse_oss_response_error(body_string)
	}
}

func CanonicalizedResourceFromObject(bucket *Bucket, object *Object) types.CanonicalizedResource {
	return types.NewCanonicalizedResource(fmt.Sprintf("/%s/%s", bucket.name, object.path))
}
