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
	return newAPIException(SERVER_ERROR, http.StatusText(http.StatusInternalServerError))
}

// 404 错误
func NotFound() *APIException {
	return newAPIException(NOT_FOUND, http.StatusText(http.StatusNotFound))
}

// 401认证错误
func UnauthorizedError() *APIException {
	return newAPIException(AUTH_ERROR, http.StatusText(http.StatusUnauthorized))
}

// 未知错误
func UnknownError(message string) *APIException {
	return newAPIException(UNKNOWN_ERROR, message)
}

// 参数错误
func ParameterError(message string) *APIException {
	return newAPIException(PARAMETER_ERROR, message)
}

// grpc错误
func GRPCError(message string) *APIException {
	return newAPIException(GRPC_ERROR, message)
}

// 其他错误
func Exception(code int, message string) *APIException {
	return newAPIException(code, message)
}
