package config

import (
	"log"
	"fmt"
	"net/http"
	"encoding/json"
	"bytes"
	"io"

	"github.com/vert3xc/bebrochka/server/internal/models"
)

func UpdateTeams(cfg *Config) error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%d/api/client/teams/", cfg.BoardHost, cfg.BoardPort), bytes.NewBuffer(nil))
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("HTTP request error:", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return err
	}
	log.Println("🔍 Raw response:", string(body))
	var teams []struct {
		Active   bool `json:"active"`
		Highlited bool `json:"highlited"`
		Id       int  `json:"id"`
		Ip       string `json:"ip"`
		Name     string `json:"name"`
	}
	if err := json.Unmarshal(body, &teams); err != nil {
		var single struct {
			Active   bool `json:"active"`
			Highlited bool `json:"highlited"`
			Id       int  `json:"id"`
			Ip       string `json:"ip"`
			Name     string `json:"name"`
		}	
		if err2 := json.Unmarshal(body, &single); err2 != nil {
			log.Println("❌ Error decoding response JSON as both array and object:", err2)
			return err
		}
		teams = append(teams, single)
	}
	for _, team := range teams {
		cfg.Teams = append(cfg.Teams, models.Team{
			Ip:   team.Ip,
			Name: team.Name,
		})
	}
	return nil	
}