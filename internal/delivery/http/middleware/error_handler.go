package middleware

import (
	"net/http"

	"backend-context-engineering-template/internal/delivery/http/dto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ErrorHandler(logger *logrus.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.WithFields(logrus.Fields{
				"error":  err,
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
			}).Error("Panic recovered")

			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_server_error",
				Message: "An internal error occurred",
			})
		}
		c.Abort()
	})
}
