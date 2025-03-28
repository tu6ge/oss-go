package oss

type EnvEmtpyError struct {
	// name string
}

// NameEmtpyError实现了 Error() 方法的对象都可以
func (e *EnvEmtpyError) Error() string {
	return "ALIYUN_KEY_ID or ALIYUN_KEY_SECRET is not empty"
}
