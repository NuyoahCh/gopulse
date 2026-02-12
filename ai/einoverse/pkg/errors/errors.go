package errors

import "fmt"

// BusinessError 业务错误类型
type BusinessError struct {
	Code    string
	Message string
}

// Error 实现 error 接口
func (e *BusinessError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// 预定义错误
var (
	ErrDocumentNotFound    = &BusinessError{Code: "DOC_NOT_FOUND", Message: "文档未找到"}
	ErrInvalidInput        = &BusinessError{Code: "INVALID_INPUT", Message: "输入参数无效"}
	ErrLLMServiceError     = &BusinessError{Code: "LLM_ERROR", Message: "LLM服务调用失败"}
	ErrApplicationNotFound = &BusinessError{Code: "APP_NOT_FOUND", Message: "申请未找到"}
)
