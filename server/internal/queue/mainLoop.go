package queue

import (
	"time"

	"github.com/vert3xc/bebrochka/server/config"
	"github.com/vert3xc/bebrochka/server/internal/models"
	"github.com/vert3xc/bebrochka/server/internal/protocols"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func submitAndUpdate(batch []models.Flag, cfg config.Config, db *gorm.DB) {
	flags := lo.Map(batch, func(f models.Flag, _ int) string { return f.Flag })
	proto := cfg.Protocol
	switch proto {
	case "ructf_http":
		submitFlags := func(flags []models.Flag, cfg config.Config) {
			protocols.RuctfHttp(flags, cfg)
		}
	}
	results := submitFlags(batch, cfg)
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

func StartLoop(cfg config.Config, db *gorm.DB) {
	batch := []models.Flag{}
    ticker := time.NewTicker(20 * time.Second)
	flags := collectFromDB(db)
	submitAndUpdate(flags, cfg, db)
    for {
        select {
        case flag := <-flagChan:
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