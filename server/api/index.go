package api

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/vert3xc/bebrochka/server/config"
)


func Index(c *gin.Context, cfg *config.Config) {
    c.HTML(http.StatusOK, "index.html", struct {
        Cfg config.Config
    }{
        Cfg: *cfg,
    })
}