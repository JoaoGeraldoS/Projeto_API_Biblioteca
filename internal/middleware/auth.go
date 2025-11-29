package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const GinContextKeyEmail = "email"
const GinContextKeyRole = "role"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokeString := c.GetHeader("Authorization")

		if tokeString == "" || len(tokeString) < 7 || strings.ToUpper(tokeString[:7]) != "BEARER " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido ou formato inválido."})
			return
		}

		tokeString = tokeString[7:]

		calims, err := VerifyToken(tokeString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Erro ao verificar token")
			return
		}

		c.Set(GinContextKeyEmail, calims.Email)
		c.Set(GinContextKeyRole, calims.Role)

		c.Next()
	}
}

func RequireRole(requiredRole string) gin.HandlerFunc {

	return func(c *gin.Context) {
		roleValue, ok := c.Get(GinContextKeyRole)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		role, ok := roleValue.(string)
		if !ok || role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Acesso negado. Requer perfil '%s'.", requiredRole),
			})
		}

		c.Next()
	}
}
