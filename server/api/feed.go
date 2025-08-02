package api

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vert3xc/bebrochka/server/config"
	"github.com/vert3xc/bebrochka/server/internal/models"
)

func Feed(c *gin.Context, cfg *config.Config) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "feed.html", gin.H{
			"feeders": cfg.Feeders,
		})
		return
	}
	for i := range cfg.Feeders {
		checkboxName := "enabled_existing_" + strconv.Itoa(i)
		cfg.Feeders[i].Enabled = c.PostForm(checkboxName) == "on"
	}
	newTeamNames := c.PostFormArray("teamname[]")
	newPeriodicities := c.PostFormArray("periodicity[]")
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid multipart form: %v", err)
		return
	}
	newFiles := form.File["pseudochecker[]"]

	if len(newTeamNames) != len(newPeriodicities) || len(newTeamNames) != len(newFiles) {
		c.String(http.StatusBadRequest, "Mismatched input lengths")
		return
	}

	for i := range newTeamNames {
		team := strings.TrimSpace(newTeamNames[i])
		periodStr := strings.TrimSpace(newPeriodicities[i])
		file := newFiles[i]

		if team == "" || periodStr == "" {
			continue
		}

		period, err := strconv.Atoi(periodStr)
		if err != nil {
			continue
		}

		filename := filepath.Base(file.Filename)
		savePath := filepath.Join("uploads", filename)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.String(http.StatusInternalServerError, "Failed to save file: %v", err)
			return
		}

		feeder := models.Feeder{
			Enabled:       true,
			Team:          team,
			Periodicity:   period,
			PseudoChecker: savePath,
			Counter:       0,
		}
		cfg.Feeders = append(cfg.Feeders, feeder)
	}

	c.Redirect(http.StatusSeeOther, "/feed")
}