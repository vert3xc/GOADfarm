package main

import (
	"log"
	"path/filepath"
	"text/template"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/vert3xc/bebrochka/server/config"
	"github.com/vert3xc/bebrochka/server/api"
	"github.com/vert3xc/bebrochka/server/internal/database"
	"github.com/vert3xc/bebrochka/server/middleware"
	"github.com/vert3xc/bebrochka/server/utils"
	"github.com/vert3xc/bebrochka/server/internal/queue"
)

func main() {
	cfg := config.Load()
	db, err := database.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	sessionSecret, apiKey, err := utils.InitSecrets()
	if err != nil {
		log.Fatalf("Failed to generate secrets: %v", err)
	}
	log.Printf("Use this API key: %s", apiKey)
	cfg.ApiToken = apiKey
	go queue.StartLoop(&cfg, db)
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"base": filepath.Base,
	})
	r.LoadHTMLGlob("/server/frontend/templates/*.html")
	r.Static("/static", "/server/frontend/static")
	r.Use(sessions.Sessions("mysession", cookie.NewStore([]byte(sessionSecret))))

	r.GET("/", middleware.AuthRequired(), func(c *gin.Context) {
		api.Index(c, &cfg)
	})
	r.GET("/feed", middleware.AuthRequired(), func(c *gin.Context) {
    	api.Feed(c, &cfg)
	})
	r.POST("/feed", middleware.AuthRequired(), func(c *gin.Context) {
    	api.Feed(c, &cfg)
	})
	r.GET("/login", func(c *gin.Context) {
    	api.Login(c, &cfg)
	})
	r.POST("/login", func(c *gin.Context) {
    	api.Login(c, &cfg)
	})

	r.GET("/api/get_config", middleware.APITokenRequired(cfg), func(c *gin.Context) {
		api.GetConfig(c, &cfg)
	})
	r.GET("/api/post_flags", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "POST only"})
	})
	r.POST("/api/post_flags", middleware.APITokenRequired(cfg), func(c *gin.Context) {
		api.PostFlags(c, db, &cfg)
	})
	r.GET("/api/list_flags", middleware.APITokenRequired(cfg), func(c *gin.Context) {
		api.ListFlags(c, db)
	})

	r.Run(fmt.Sprintf(":%d", cfg.ServerPort))
}