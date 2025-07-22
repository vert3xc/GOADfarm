package main

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/vert3xc/bebrochka/server/config"
	"github.com/vert3xc/bebrochka/server/api"
	"github.com/vert3xc/bebrochka/server/internal/database"
	"github.com/vert3xc/bebrochka/server/internal/utils"
)

func main() {
	cfg := config.LoadConfig()
	db, err := database.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	r := gin.Default()
	r.Use(sessions.Sessions("mysession", cookie.NewStore([]byte(utils.GenerateSecretKey()))))

	r.GET("/api/config", func(c *gin.Context) {
		api.GetConfig(c, cfg)
	})

	r.GET("/api/post_flags", func(c *gin.Context) {
		api.GetConfig(c, db)
	})

	r.Run(":5000")
}