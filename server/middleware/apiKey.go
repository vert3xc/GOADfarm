package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/vert3xc/bebrochka/server/config"
)

func APITokenRequired(cfg config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("X-API-Token")
        if token != cfg.ApiToken {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API token"})
            return
        }
        c.Next()
    }
}