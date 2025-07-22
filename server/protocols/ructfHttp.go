package protocols

import (
    "net/http"
    "bytes"
	"fmt"
	"encoding/json"
	"strings"

	"github.com/vert3xc/bebrochka/server/config"
	"github.com/vert3xc/bebrochka/server/internal/models"
)

func RuctfHttp(flags []string, cfg config.Config) {
	jsonData, err := json.Marshal(flags)
	if err != nil {
		fmt.Println("Error marshaling flags:", err)
		return
	}
    client := &http.Client{}
    req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s:%d/flags", cfg.BoardHost, cfg.BoardPort), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("X-Team-Token", cfg.TeamToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP request error:", err)
		return
	}
	defer resp.Body.Close()
	var responses []struct {
        Msg  string `json:"msg"`
        Flag string `json:"flag"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&responses); err != nil {
        fmt.Println("Error decoding response JSON:", err)
        return
    }
	var results []models.SubmitResult
    unknownResponses := make(map[string]bool)

    for _, item := range responses {
        msg := strings.TrimSpace(item.Msg)
        cleaned := strings.ReplaceAll(msg, fmt.Sprintf("[%s] ", item.Flag), "")
        cleanedLower := strings.ToLower(cleaned)

        foundStatus := models.StatusQueued
        for status, substrings := range models.ResponsesMap {
            for _, substr := range substrings {
                if strings.Contains(cleanedLower, substr) {
                    foundStatus = status
                    break
                }
            }
            if foundStatus != models.StatusQueued {
                break
            }
        }

        if foundStatus == models.StatusQueued && !unknownResponses[cleaned] {
            unknownResponses[cleaned] = true
            fmt.Printf("⚠️ Unknown checksystem response (flag will be resent): %s\n", cleaned)
        }

        results = append(results, models.SubmitResult{
            Flag:   item.Flag,
            Status: foundStatus,
            Msg:    cleaned,
        })
    }

    return results
}