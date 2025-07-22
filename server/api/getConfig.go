package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vert3xc/bebrochka/server/config"
)

type PublicConfig struct {
    Teams           []string `json:"teams"`
    FlagFormat      string   `json:"flag_format"`
    TeamToken       string   `json:"team_token"`
    Protocol        string   `json:"protocol"`
    BoardHost       string   `json:"board_host"`
    BoardPort       int      `json:"board_port"`
    FlagLimit       int      `json:"flag_limit"`
    FlagSubmitPeriod int     `json:"flag_submit_period"`
    FlagLifeTime    int      `json:"flag_life_time"`
}

func GetConfig(c *gin.Context, cfg config.Config) {
	c.JSON(http.StatusOK, PublicConfig{
		Teams:           cfg.Teams,
		FlagFormat:      cfg.FlagFormat,
		TeamToken:       cfg.TeamToken,
		Protocol:        cfg.Protocol,
		BoardHost:       cfg.BoardHost,
		BoardPort:       cfg.BoardPort,
		FlagLimit:       cfg.FlagLimit,
		FlagSubmitPeriod: cfg.FlagSubmitPeriod,
		FlagLifeTime:    cfg.FlagLifeTime,
	})
}