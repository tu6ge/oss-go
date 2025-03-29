package oss

import "fmt"

type EnvEmtpyError struct {
	// name string
}

// NameEmtpyError实现了 Error() 方法的对象都可以
func (e *EnvEmtpyError) Error() string {
	return "ALIYUN_KEY_ID or ALIYUN_KEY_SECRET is not empty"
}

type OssResponseError struct {
	Code         string
	Message      string
	RequestId    string
	RecommendDoc string
}

func parse_oss_response_error(xml string) *OssResponseError {
	code := parse_item(xml, "Code")
	message := parse_item(xml, "Message")
	request_id := parse_item(xml, "RequestId")
	recommend_doc := parse_item(xml, "RecommendDoc")

	return &OssResponseError{code, message, request_id, recommend_doc}
}

func (e *OssResponseError) Error() string {
	return fmt.Sprintf("oss return: %s", e.Message)
}
