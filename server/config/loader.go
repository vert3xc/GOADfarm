package config

import (
	"log"
	"os"
	"fmt"
)

type Config struct {
	Teams        []string
	FlagFormat   string
	TeamToken    string
	Protocol     string
	BoardHost    string
	BoardPort    int
	FlagLimit    int
	FlagSubmitPeriod   int
	FlagLifeTime int
	Password     string
	ApiToken     string
	DSN          string
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func Load() Config {
	teams := []string{}
	ipFormat := getEnv("IP_FORMAT", "10.0.0.%d")
	lowerTeamBound := strconv.Atoi(getEnv("LOWER_BOUND", 1))
	upperTeamBound := strconv.Atoi(getEnv("UPPER_BOUND", 30))
	for i := lowerTeamBound; i < upperTeamBound; i++ {
		teams = append(teams, fmt.Sprintf(ipFormat, i))
	}
	return Config{
		Teams:      teams,
		FlagFormat: getEnv("FLAG_FORMAT", "[A-Z0-9]{31}="),
		Protocol:   getEnv("PROTOCOL", "ructf_http"),
		BoardHost:  getEnv("BOARD_HOST", "127.0.0.1"),
		BoardPort:  getEnv("BOARD_PORT", 8080),
		FlagLimit:  getEnv("FLAG_LIMIT", 50),
		FlagSubmitPeriod: getEnv("FLAG_SUBMIT_PERIOD", 5),
		FlagLifeTime: getEnv("FLAG_LIFETIME", 300),
		Password:   getEnv("PASSWORD", "bebra"),
		ApiToken:   getEnv("API_TOKEN", "token"),
		TeamToken:  getEnv("TEAM_TOKEN", "token"),
		DSN:       getEnv("POSTGRES_DSN", "host=localhost user=postgres password=password dbname=bebra port=5432 sslmode=disable"),
	}
}
