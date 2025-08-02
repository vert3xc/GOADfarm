package queue

import (
	"time"
    "fmt"

	"github.com/vert3xc/bebrochka/server/config"
	"github.com/vert3xc/bebrochka/server/internal/models"
	"github.com/vert3xc/bebrochka/server/protocols"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func submitAndUpdate(batch []models.Flag, cfg *config.Config, db *gorm.DB) {
	flags := lo.Map(batch, func(f models.Flag, _ int) string { return f.Flag })
	proto := cfg.Protocol
    submitFlags := func(flags []string, cfg *config.Config) ([]models.SubmitResult, error) {
            return nil, nil
        }
	switch proto {
	case "ructf_http":
		submitFlags = func(flags []string, cfg *config.Config) ([]models.SubmitResult, error) {
            return protocols.RuctfHttp(flags, cfg)
        }
	}
	results, err := submitFlags(flags, cfg)
	if err != nil {
		fmt.Println("Error submitting flags:", err)
		return
	}
	for _, res := range results {
        db.Model(&models.Flag{}).
            Where("flag = ?", res.Flag).
            Updates(map[string]interface{}{
                "status":   res.Status,
                "response": res.Msg,
            })
    }
}

func collectFromDB(db *gorm.DB) []models.Flag {
	var flags []models.Flag
	db.Where("status = ?", models.StatusQueued).Find(&flags)
	return flags
}

func StartLoop(cfg *config.Config, db *gorm.DB) {
	batch := []models.Flag{}
    ticker := time.NewTicker(20 * time.Second)
	flags := collectFromDB(db)
    if len(flags) > 0 {
	    submitAndUpdate(flags, cfg, db)
    }
    for {
        select {
        case flag := <-FlagChan:
            batch = append(batch, flag)
            if len(batch) >= cfg.FlagLimit {
                submitAndUpdate(batch, cfg, db)
                batch = nil
            }
        case <-ticker.C:
            if len(batch) > 0 {
                submitAndUpdate(batch, cfg, db)
                batch = nil
            }
        }
    }
}