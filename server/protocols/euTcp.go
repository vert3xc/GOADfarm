package protocols

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"encoding/json"
	"net/http"
	"io"
	"time"

	"github.com/vert3xc/bebrochka/server/config"
	"github.com/vert3xc/bebrochka/server/internal/models"
	"github.com/vert3xc/bebrochka/server/utils"
)

func UpdateTeamsEuctf(cfg *config.Config) error {
    var url string
    if cfg.BoardPort != -1 {
        url = fmt.Sprintf("http://%s:%d/competition/teams.json", cfg.BoardHost, cfg.BoardPort)
    } else {
        url = fmt.Sprintf("http://%s/competition/teams.json", cfg.BoardHost)
    }

    client := &http.Client{}
    req, err := http.NewRequest(http.MethodGet, url, nil)
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

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error reading response body:", err)
        return err
    }

    var apiResp struct {
        Teams []int `json:"teams"`
    }
    if err := json.Unmarshal(body, &apiResp); err != nil {
        log.Println("Error decoding JSON:", err)
        return err
    }

    cfg.Teams = cfg.Teams[:0]

    for _, team := range apiResp.Teams {
        cfg.Teams = append(cfg.Teams, models.Team{
            Ip:   fmt.Sprintf("10.66.%d.123", team),
            Name: "Marcusov",
        })
    }

    return nil
}

func EuTCP(flags []string, cfg *config.Config) ([]models.SubmitResult, error) {
	var servAddr string
	if cfg.BoardPort != -1 {
		servAddr = fmt.Sprintf("%s:%d", cfg.BoardHost, cfg.BoardPort)
	} else {
		servAddr = cfg.BoardHost
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		log.Println("ResolveTCPAddr failed:", err)
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println("Dial failed:", err)
		return nil, err
	}
	defer conn.Close()
	timeout := 5 * time.Second
	conn.SetDeadline(time.Now().Add(timeout))
	_, err = utils.RecvUntil(conn, "\n\n")
	if err != nil {
		log.Println("Error receiving welcome message:", err)
		return nil, err
	}
	toSend := strings.Join(flags, "\n") + "\n"
	if err := utils.SendLine(conn, toSend); err != nil {
		log.Println("Error sending flags:", err)
		return nil, err
	}
	results := make([]models.SubmitResult, 0, len(flags))
	reader := bufio.NewReader(conn)
	for i := 0; i < len(flags); i++ {
		conn.SetReadDeadline(time.Now().Add(timeout))
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading response:", err)
			break
		}
		parts := strings.Fields(strings.TrimSpace(response))
		if len(parts) == 0 {
			log.Println("Empty response received")
			continue
		}
		flag := parts[0]
		msg := ""
		status := models.StatusQueued
		if len(parts) == 2 {
			switch parts[0] {
			case "OK":
				msg = "Flag accepted"
				status = models.StatusAccepted
			case "DUP":
				msg = "Flag already submitted"
				status = models.StatusRejected
			case "OWN":
				msg = "Flag belongs to your team"
				status = models.StatusRejected
			case "INV":
				msg = "Invalid flag"
				status = models.StatusRejected
			case "ERR":
				msg = "The server encountered an internal error"
				status = models.StatusRejected
			default:
				msg = parts[1]
			}
		} else if len(parts) >= 3 {
			cleanedLower := strings.ToLower(parts[1])
			msg = parts[2]
			for s, substrings := range models.ResponsesMap {
				for _, substr := range substrings {
					if strings.Contains(cleanedLower, substr) {
						status = s
						break
					}
				}
				if status != models.StatusQueued {
					break
				}
			}
		}
		results = append(results, models.SubmitResult{
			Flag:   flag,
			Status: status,
			Msg:    msg,
		})
	}
	return results, nil
}
