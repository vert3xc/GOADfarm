package config

import (
	"os"
	"fmt"
	"strconv"
	"log"
	"sync"

	"github.com/vert3xc/bebrochka/server/internal/models"
)

type Config struct {
	Teams        []models.Team
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
	Feeders      []models.Feeder
	Mut          sync.Mutex
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func Load() Config {
	teams := []models.Team{}
	ipFormat := getEnv("IP_FORMAT", "10.0.0.%d")
	lowerTeamBound, err := strconv.Atoi(getEnv("LOWER_BOUND", "1"))
	if err != nil {
		log.Fatalf("Invalid LOWER_BOUND: %v", err)
	}
	upperTeamBound, err := strconv.Atoi(getEnv("UPPER_BOUND", "30"))
	if err != nil {
		log.Fatalf("Invalid UPPER_BOUND: %v", err)
	}
	for i := lowerTeamBound; i < upperTeamBound; i++ {
		teams = append(teams, models.Team{
			Ip:   fmt.Sprintf(ipFormat, i),
			Name: fmt.Sprintf("Team #%d", i),
		})
	}
	boardPort, err := strconv.Atoi(getEnv("BOARD_PORT", "8080"))
	if err != nil {
		log.Fatalf("Invalid BOARD_PORT: %v", err)
	}
	flagLimit, err := strconv.Atoi(getEnv("FLAG_LIMIT", "50"))
	if err != nil {
		log.Fatalf("Invalid FLAG_LIMIT: %v", err)
	}
	flagSubmitPeriod, err := strconv.Atoi(getEnv("FLAG_SUBMIT_PERIOD", "5"))
	if err != nil {
		log.Fatalf("Invalid FLAG_SUBMIT_PERIOD: %v", err)
	}
	flagLifeTime, err := strconv.Atoi(getEnv("FLAG_LIFETIME", "300"))
	if err != nil {
		log.Fatalf("Invalid FLAG_LIFETIME: %v", err)
	}
	return Config{
		Teams:      teams,
		FlagFormat: getEnv("FLAG_FORMAT", "[A-Z0-9]{31}="),
		Protocol:   getEnv("PROTOCOL", "ructf_http"),
		BoardHost:  getEnv("BOARD_HOST", "127.0.0.1"),
		BoardPort:  boardPort,
		FlagLimit:  flagLimit,
		FlagSubmitPeriod: flagSubmitPeriod,
		FlagLifeTime: flagLifeTime,
		Password:   getEnv("PASSWORD", "bebra"),
		ApiToken:   "apitoken",
		TeamToken:  getEnv("TEAM_TOKEN", "token"),
		DSN:       getEnv("POSTGRES_DSN", "host=db user=postgres password=password dbname=bebra port=5432 sslmode=disable"),
		Feeders: []models.Feeder{},
		Mut:      sync.Mutex{},
	}
}
