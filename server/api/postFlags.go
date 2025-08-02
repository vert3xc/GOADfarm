package api

import (
	"net/http"
	"time"
	"os/exec"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"github.com/vert3xc/bebrochka/server/config"
	"github.com/vert3xc/bebrochka/server/utils"
	"github.com/vert3xc/bebrochka/server/internal/models"
	"github.com/vert3xc/bebrochka/server/internal/queue"
	"github.com/samber/lo"
)

type FlagInput struct {
	Flag   string `json:"flag"`
	Sploit string `json:"sploit"`
	Team   string `json:"team"`
}

func processFeeding(flag FlagInput, cfg *config.Config) {
	idx, ok := utils.Contains(lo.Map(cfg.Feeders, func(f models.Feeder, _ int) string {
		return f.Team
	}), flag.Team)
	if !ok {
		return
	}
	if !cfg.Feeders[idx].Enabled {
		return
	}
	cfg.Mut.Lock()
	shouldRun := cfg.Feeders[idx].Counter%cfg.Feeders[idx].Periodicity == 0
	cfg.Feeders[idx].Counter++
	cfg.Mut.Unlock()
	if shouldRun {
		log.Printf("Running pseudochecker for flag: %s\n", flag.Flag)
		cmd := exec.Command("python", cfg.Feeders[idx].PseudoChecker, flag.Flag)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error occurred while executing pseudochecker: %v\nOutput: %s\n", err, output)
		}
		log.Printf("Pseudochecker output: %s\n", output)
	}
}

func PostFlags(c *gin.Context, db *gorm.DB, cfg *config.Config) {
	var flags []FlagInput
	if err := c.ShouldBindJSON(&flags); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	var toInsert []models.Flag
	curTime := time.Now()
	for _, item := range flags {
		go processFeeding(item, cfg)
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