package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vert3xc/bebrochka/server/config"
	"github.com/vert3xc/bebrochka/server/internal/models"
)

type PublicConfig struct {
    Teams           []models.Team `json:"TEAMS"`
    FlagFormat      string         `json:"FLAG_FORMAT"`
    Protocol        string         `json:"SYSTEM_PROTOCOL"`
    BoardHost       string         `json:"SYSTEM_HOST"`
    BoardPort       int            `json:"SYSTEM_PORT"`
    FlagLimit       int            `json:"SUBMIT_FLAG_LIMIT"`
    FlagSubmitPeriod int           `json:"SUBMIT_PERIOD"`
    FlagLifeTime    int            `json:"FLAG_LIFETIME"`
}

func GetConfig(c *gin.Context, cfg *config.Config) {
	c.JSON(http.StatusOK, PublicConfig{
		Teams:           cfg.Teams,
		FlagFormat:      cfg.FlagFormat,
		Protocol:        cfg.Protocol,
		BoardHost:       cfg.BoardHost,
		BoardPort:       cfg.BoardPort,
		FlagLimit:       cfg.FlagLimit,
		FlagSubmitPeriod: cfg.FlagSubmitPeriod,
		FlagLifeTime:    cfg.FlagLifeTime,
	})
}