package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func NewApiError(status int, code string, message string, err error) *APIError {
	return &APIError{
		Status:  status,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *APIError) Messager(msg string) *APIError {
	newErr := *e
	newErr.Message = msg
	return &newErr
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("API Error: Status %d, Code %s, Message %s, Internal Err: %v", e.Status, e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("API Error: Status %d, Code %s, Message %s", e.Status, e.Code, e.Message)
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			reqPath := c.Request.URL.Path

			if apiErr, ok := err.(*APIError); ok {
				c.AbortWithStatusJSON(apiErr.Status, gin.H{
					"status":  apiErr.Status,
					"code":    apiErr.Code,
					"message": apiErr.Message,
					"path":    reqPath,
				})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"code":    "INTERNAL_SERVER_ERROR",
					"message": "Ocorreu um erro inesperado",
					"path":    reqPath,
				})
			}
			return
		}
	}
}

var (
	NotFound    = NewApiError(http.StatusNotFound, "NOT_FOUND", "Recurno não encontrado", nil)
	BadRequest  = NewApiError(http.StatusBadRequest, "BAD_REQUEST", "Solicitação invalida", nil)
	InternalErr = NewApiError(http.StatusInternalServerError, "INTERNAL_ERROR", "Erro interno ocorrido", nil)
)
