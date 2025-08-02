package api

import (
	"net/http"

	"github.com/vert3xc/bebrochka/server/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
)

type LoginInput struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func Login(c *gin.Context, cfg *config.Config) {
	if c.Request.Method == http.MethodGet {
        c.HTML(http.StatusOK, "login.html", gin.H{})
        return
    }

    var input LoginInput
    if err := c.ShouldBind(&input); err != nil {
        c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid input"})
        return
    }

    if input.Username == "bebroid" && input.Password == cfg.Password {
        session := sessions.Default(c)
        session.Set("user", input.Username)
        session.Save()
        c.Redirect(http.StatusFound, "/")
        return
    }

    c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid credentials"})
}