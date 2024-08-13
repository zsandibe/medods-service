package v1

import (
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, status string, code int, err error) {
	c.AbortWithStatusJSON(code, errorResponse{Status: status, Code: code, Message: err.Error()})
}
