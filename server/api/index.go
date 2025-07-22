package api

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/vert3xc/bebrochka/server/internal/models"
    "gorm.io/gorm"
)

type IndexPageData struct {
    DistinctSploit           []string
    DistinctTeam             []string
    DistinctStatus           []string
    SelectedSploit           string
    SelectedTeam             string
    SelectedFlag             string
    SelectedTimeSince        string
    SelectedTimeUntil        string
    SelectedStatus           string
    SelectedChecksystemResponse string
    ServerTZName             string
    FlagFormat               string
    Flags                    []models.Flag
    TotalFlags               int
    Pages                    []int
    CurrentPage              int
    QueryStringExceptPage    string
}

func Index(c *gin.Context, db *gorm.DB, flagFormat string, serverTZName string) {
    selectedSploit := c.Query("sploit")
    selectedTeam := c.Query("team")
    selectedFlag := c.Query("flag")
    selectedTimeSince := c.Query("time-since")
    selectedTimeUntil := c.Query("time-until")
    selectedStatus := c.Query("status")
    selectedChecksystemResponse := c.Query("checksystem_response")
    pageStr := c.DefaultQuery("page", "1")
    page, _ := strconv.Atoi(pageStr)
    if page < 1 {
        page = 1
    }
    pageSize := 50

    var flags []models.Flag
    query := db.Model(&models.Flag{})
    if selectedSploit != "" {
        query = query.Where("sploit = ?", selectedSploit)
    }
    if selectedTeam != "" {
        query = query.Where("team = ?", selectedTeam)
    }
    if selectedFlag != "" {
        query = query.Where("flag ILIKE ?", "%"+selectedFlag+"%")
    }
    if selectedStatus != "" {
        query = query.Where("status = ?", selectedStatus)
    }
    if selectedChecksystemResponse != "" {
        query = query.Where("response ILIKE ?", "%"+selectedChecksystemResponse+"%")
    }

    var total int64
    query.Count(&total)
    query = query.Order("time DESC").Offset((page - 1) * pageSize).Limit(pageSize)
    query.Find(&flags)

    var distinctSploit, distinctTeam, distinctStatus []string
    db.Model(&models.Flag{}).Distinct().Pluck("sploit", &distinctSploit)
    db.Model(&models.Flag{}).Distinct().Pluck("team", &distinctTeam)
    db.Model(&models.Flag{}).Distinct().Pluck("status", &distinctStatus)

    totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
    var pages []int
    for i := 1; i <= totalPages; i++ {
        pages = append(pages, i)
    }

    queryStringExceptPage := ""
    for key, vals := range c.Request.URL.Query() {
        if key == "page" {
            continue
        }
        for _, val := range vals {
            queryStringExceptPage += "&" + key + "=" + val
        }
    }

    data := IndexPageData{
        DistinctSploit:           distinctSploit,
        DistinctTeam:             distinctTeam,
        DistinctStatus:           distinctStatus,
        SelectedSploit:           selectedSploit,
        SelectedTeam:             selectedTeam,
        SelectedFlag:             selectedFlag,
        SelectedTimeSince:        selectedTimeSince,
        SelectedTimeUntil:        selectedTimeUntil,
        SelectedStatus:           selectedStatus,
        SelectedChecksystemResponse: selectedChecksystemResponse,
        ServerTZName:             serverTZName,
        FlagFormat:               flagFormat,
        Flags:                    flags,
        TotalFlags:               int(total),
        Pages:                    pages,
        CurrentPage:              page,
        QueryStringExceptPage:    queryStringExceptPage,
    }

    c.HTML(http.StatusOK, "index.html", data)
}