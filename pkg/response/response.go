package httpresponse

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status     bool        `json:"status"`
	StatusCode int32       `json:"statusCode"`
	Message    *string     `json:"message,omitempty"`
	Error      *string     `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func NewResponse() *Response {
	return &Response{
		Status:     true,
		StatusCode: 200,
	}
}

func Success(data interface{}, message string) *Response {
	return &Response{
		Status:     true,
		StatusCode: 200,
		Message:    &message,
		Data:       data,
	}
}

func Error(statusCode int32, errorMessage string) *Response {
	return &Response{
		Status:     false,
		StatusCode: statusCode,
		Error:      &errorMessage,
	}
}

func (r *Response) WithData(data interface{}) *Response {
	r.Data = data
	return r
}

func (r *Response) WithMessage(message string) *Response {
	r.Message = &message
	return r
}

func (r *Response) WithError(err string) *Response {
	r.Error = &err
	r.Status = false
	return r
}

func (r *Response) WithStatusCode(code int32) *Response {
	r.StatusCode = code
	return r
}

type Responses struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Status:  statusCode,
		Message: message,
	})
}

func SendErrorResponseWithDetails(c *gin.Context, statusCode int, message, details string) {
	c.JSON(statusCode, ErrorResponse{
		Status:  statusCode,
		Message: message,
		Details: details,
	})
}

func SendSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}

func ResponseTimeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		log.Printf("Response time for %s %s: %v\n", c.Request.Method, c.Request.URL.Path, duration)
	}
}
