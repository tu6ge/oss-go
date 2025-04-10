package oss

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tu6ge/oss-go/types"
)

type PartsUpload struct {
	path      string
	upload_id string
	etag_list []etag_struct
}

type etag_struct struct {
	index   int
	content string
}

func NewPartsUpload(path string) PartsUpload {
	return PartsUpload{path, "", []etag_struct{}}
}

func (m PartsUpload) ToUrl(bucket *Bucket) url.URL {
	url := bucket.ToUrl()
	url.Path = m.path
	return url
}

func (m PartsUpload) InitMulit(client *Client) error {
	bucket := client.Bucket
	url := m.ToUrl(&bucket)
	method := "POST"

	resource := canonicalized_resource(&bucket, &m)
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

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	body_string := string(body)
	m.upload_id = parse_upload_id(body_string)
	if len(m.upload_id) == 0 {
		return errors.New("not found upload_id")
	}

	return nil
}

func (m PartsUpload) UploadPart(index int, con []byte, client *Client) error {
	bucket := client.Bucket
	url := m.ToUrl(&bucket)
	method := "PUT"

	resource := canonicalized_resource_part(&bucket, &m, index, m.upload_id)

	headers := map[string]string{
		"Content-Length": strconv.Itoa(len(con)),
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
	etag := resp.Header.Get("ETag")
	if len(etag) == 0 {
		return errors.New("not found etag header")
	}

	m.etag_list = append(m.etag_list, etag_struct{index, etag})

	return nil
}

func (m PartsUpload) etag_list_xml() string {
	list := ""
	for _, item := range m.etag_list {
		list += fmt.Sprintf("<Part><PartNumber>%d</PartNumber><ETag>%s</ETag></Part>", item.index, item.content)
	}

	return fmt.Sprintf("<CompleteMultipartUpload>%s</CompleteMultipartUpload>", list)
}

func (m PartsUpload) Complete(client *Client) error {
	bucket := client.Bucket
	url := m.ToUrl(&bucket)
	method := "POST"

	resource := canonicalized_resource_complete(&bucket, &m, m.upload_id)

	xml := m.etag_list_xml()

	headers := map[string]string{
		"Content-Length": strconv.Itoa(len(xml)),
	}
	headers = client.AuthorizationHeader(method, resource, headers)

	req, err := http.NewRequest(method, url.String(), bytes.NewReader([]byte(xml)))
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
	if http_status_ok(resp.StatusCode) {
		return nil
	}
	return errors.New("complete failed")
}

func canonicalized_resource(bucket *Bucket, object *PartsUpload) types.CanonicalizedResource {
	return types.NewCanonicalizedResource(fmt.Sprintf("/%s/%s?upload", bucket.name, object.path))
}

func canonicalized_resource_part(bucket *Bucket, object *PartsUpload, index int, upload_id string) types.CanonicalizedResource {
	return types.NewCanonicalizedResource(fmt.Sprintf("/%s/%s?partNumber=%d&uploadId=%s", bucket.name, object.path, index, upload_id))
}

func canonicalized_resource_complete(bucket *Bucket, object *PartsUpload, upload_id string) types.CanonicalizedResource {
	return types.NewCanonicalizedResource(fmt.Sprintf("/%s/%s?uploadId=%s", bucket.name, object.path, upload_id))
}

func parse_upload_id(xml string) string {
	start := strings.Index(xml, "<UploadId>")

	if start == -1 {
		return ""
	}
	end := strings.Index(xml, "</UploadId>")
	if end == -1 {
		return ""
	}
	return xml[start+10 : end]
}
