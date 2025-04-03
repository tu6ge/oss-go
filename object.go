package oss

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/tu6ge/oss-go/types"
)

type Objects struct {
	List      []Object
	NextToken string
	query     types.ObjectQuery
}

func (objs Objects) NextList(client *Client) (Objects, error) {
	if len(objs.NextToken) == 0 {
		return Objects{}, &NoFoundMoreObject{}
	}
	objs.query.Insert(types.QUERY_CONTINUATION_TOKEN, objs.NextToken)
	return client.Bucket.ObjectQuery(objs.query).GetObjects(client)
}

type NoFoundMoreObject struct{}

func (n *NoFoundMoreObject) Error() string {
	return "no found more object"
}

type Object struct {
	path         string
	content      []byte
	content_type string
	copy_source  string
	errors       error
}

func (obj Object) String() string {
	return obj.path
}

func NewObject(path string) Object {
	return Object{path, nil, "", "", nil}
}

func (obj Object) ToUrl(bucket *Bucket) url.URL {
	url := bucket.ToUrl()
	url.Path = obj.path
	return url
}

func (obj Object) Content(con []byte) Object {
	obj.content = con
	return obj
}

func (obj Object) File(reader io.Reader) Object {
	con, err := io.ReadAll(reader)
	if err != nil {
		obj.errors = err
		return obj
	}
	obj.content = con
	return obj
}

func (obj Object) FilePath(name string) Object {
	f, err := os.Open(name)
	if err != nil {
		obj.errors = err
		return obj
	}
	defer f.Close()

	return obj.File(f)
}

func (obj Object) ContentType(con string) Object {
	obj.content_type = con
	return obj
}

func (obj Object) Upload(client *Client) error {
	if obj.errors != nil {
		return obj.errors
	}

	bucket := client.Bucket
	url := obj.ToUrl(&bucket)
	method := "PUT"

	resource := CanonicalizedResourceFromObject(&bucket, &obj)
	headers := make(map[string]string)
	if len(obj.content_type) > 0 {
		headers["Content-Type"] = obj.content_type
	}
	headers = client.AuthorizationHeader(method, resource, headers)

	if len(obj.content) == 0 {
		headers["Content-Length"] = "0"
	}

	req, err := http.NewRequest(method, url.String(), bytes.NewReader(obj.content))
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
	bucket := client.Bucket
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

func (obj Object) CopySource(source string) Object {
	obj.copy_source = source
	return obj
}

func (obj Object) Copy(client *Client) error {
	bucket := client.Bucket
	url := obj.ToUrl(&bucket)
	method := "PUT"

	resource := CanonicalizedResourceFromObject(&bucket, &obj)
	headers := make(map[string]string)
	if len(obj.copy_source) == 0 {
		return errors.New("not found copy source")
	}
	headers["x-oss-copy-source"] = obj.copy_source
	if len(obj.content_type) > 0 {
		headers["Content-Type"] = obj.content_type
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
	bucket := client.Bucket
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
