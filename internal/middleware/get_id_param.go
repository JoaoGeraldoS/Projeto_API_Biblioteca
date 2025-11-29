package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetIdParam(c *gin.Context) (int64, error) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, BadRequest.Messager("Id invalido")
	}
	return id, nil
}
