package api

import (
	"net/http"
	"strconv"
	"math"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/vert3xc/bebrochka/server/internal/models"
)

func ListFlags(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit
	query := db.Model(&models.Flag{}).Order("time desc")

	if sploit := c.Query("sploit"); sploit != "" {
		query = query.Where("sploit ILIKE ?", "%"+sploit+"%")
	}
	if team := c.Query("team"); team != "" {
		query = query.Where("team ILIKE ?", "%"+team+"%")
	}
	if flag := c.Query("flag"); flag != "" {
		query = query.Where("flag ILIKE ?", "%"+flag+"%")
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if response := c.Query("response"); response != "" {
		query = query.Where("response ILIKE ?", "%"+response+"%")
	}

	var total int64
	query.Count(&total)

	var flags []models.Flag
	query.Offset(offset).Limit(limit).Find(&flags)

	c.JSON(http.StatusOK, gin.H{
		"flags":       flags,
		"total_pages": int(math.Ceil(float64(total) / float64(limit))),
	})
}