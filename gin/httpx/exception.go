package httpx

import "net/http"

type APIException struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 实现接口
func (e *APIException) Error() string {
	return e.Message
}

func newAPIException(code int, msg string) *APIException {
	return &APIException{
		Code:    code,
		Message: msg,
	}
}

// 500 错误处理
func ServerError() *APIException {
	return newAPIException(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

// 404 错误
func NotFound() *APIException {
	return newAPIException(http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

// 未知错误
func UnknownError(message string) *APIException {
	return newAPIException(http.StatusForbidden, message)
}

// 参数错误
func ParameterError(message string) *APIException {
	return newAPIException(http.StatusBadRequest, message)
}

// grpc错误
func GRPCError(message string) *APIException {
	return newAPIException(http.StatusInternalServerError, message)
}
