package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vert3xc/bebrochka/server/config"
	"github.com/vert3xc/bebrochka/server/internal/database"
	"github.com/vert3xc/bebrochka/server/internal/models"
	"github.com/vert3xc/bebrochka/server/internal/queue"
)

type FlagInput struct {
	Flag   string `json:"flag"`
	Sploit string `json:"sploit"`
	Team   string `json:"team"`
}

func PostFlags(c *gin.Context, db *gorm.DB) {
	var flags []FlagInput
	if err := c.ShouldBindJSON(&flags); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	var toInsert []models.Flag
	curTime := time.Now()
	for _, item := range flags {
		toInsert = append(toInsert, models.Flag{
			Flag:   item.Flag,
			Sploit: item.Sploit,
			Team:   item.Team,
			Time:   curTime,
			Status: models.StatusQueued,
		})
	}

	if len(toInsert) > 0 {
		for _, flag := range toInsert {
			db.Clauses(clause.OnConflict{DoNothing: true}).Create(&flag)
			queue.FlagChan <- flag
		}
	}

	c.String(http.StatusOK, "")
}