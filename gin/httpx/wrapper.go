package httpx

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

type HandlerFunc func(c *gin.Context) error

func Wrapper(handler HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := handler(c)
		if err != nil {
			var apiException *APIException
			if h, ok := err.(*APIException); ok {
				apiException = h
			} else if e, ok := status.FromError(err); ok {
				apiException = GRPCError(e.Message())
			} else if e, ok := err.(error); ok {
				apiException = UnknownError(e.Error())
			} else {
				apiException = ServerError()
			}
			// apiException.Request = c.Request.Method + " " + c.Request.URL.String()
			c.JSON(apiException.Code, apiException)
			return
		}
	}
}
