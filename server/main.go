package main

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/vert3xc/bebrochka/server/config"
	"github.com/vert3xc/bebrochka/server/api"
	"github.com/vert3xc/bebrochka/server/internal/database"
	"github.com/vert3xc/bebrochka/server/internal/middleware"
	"github.com/vert3xc/bebrochka/server/internal/utils"
	"github.com/vert3xc/bebrochka/server/internal/queue"
)

func main() {
	cfg := config.LoadConfig()
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
	go queue.StartLoop(cfg, db)
	r := gin.Default()
	r.LoadHTMLGlob("server/frontend/templates/*.html")
	r.Static("/static", "./frontend/static")
	r.Use(sessions.Sessions("mysession", cookie.NewStore([]byte(sessionSecret))))

	r.GET("/", middleware.AuthRequired(), func(c *gin.Context) {
		api.Index(c, db, cfg.FlagFormat, "UTC")
	})
	r.GET("/login", func(c *gin.Context) {
    	api.Login(c, cfg)
	})
	r.POST("/login", func(c *gin.Context) {
    	api.Login(c, cfg)
	})

	r.GET("/api/config", middleware.APITokenRequired(cfg), func(c *gin.Context) {
		api.GetConfig(c, cfg)
	})
	r.POST("/api/post_flags", middleware.APITokenRequired(cfg), func(c *gin.Context) {
		api.PostFlags(c, db)
	})

	r.Run(":5001")
}