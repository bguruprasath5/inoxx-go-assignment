package response

import "github.com/gin-gonic/gin"

// Response struct
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorJSON : json error response function
func ErrorJSON(c *gin.Context, statusCode int, err string) {
	c.AbortWithStatusJSON(statusCode, &Response{
		Success: false,
		Message: err,
		Data:    nil,
	})
}

// SuccessJSON : json error response function
func SuccessJSON(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, &Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}
