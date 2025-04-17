package oss

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/tu6ge/oss-go/types"
)

type PartsDownload struct {
	path      string
	part_size int
	file_path string
}

func NewPartsDownload(path string) PartsDownload {
	return PartsDownload{path, 1024 * 1024, ""}
}

func (obj PartsDownload) ToUrl(bucket *Bucket) url.URL {
	url := bucket.ToUrl()
	url.Path = obj.path
	return url
}

func (p PartsDownload) PartSize(size int) PartsDownload {
	p.part_size = size
	return p
}

func (p PartsDownload) FilePath(path string) PartsDownload {
	p.file_path = path
	return p
}

func (p PartsDownload) Download(client *Client) error {
	bucket := client.Bucket
	url := p.ToUrl(&bucket)
	method := "GET"

	resource := types.NewCanonicalizedResource(fmt.Sprintf("/%s/%s", bucket.name, p.path))

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

	// 检查 HTTP 状态码
	if !http_status_ok(resp.StatusCode) {
		// 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		body_string := string(body)
		return parse_oss_response_error(body_string)
	}

	// 创建本地文件用于保存内容
	outFile, err := os.Create(p.file_path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 设置缓冲大小
	bufferSize := p.part_size
	writer := bufio.NewWriterSize(outFile, bufferSize)

	// 使用 io.CopyBuffer 实现分片复制
	buf := make([]byte, bufferSize)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := writer.Write(buf[:n]); writeErr != nil {
				return writeErr
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	// 刷新缓冲区
	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}
