package models

import (
	"gorm.io/gorm"
	"time"
)

type FlagStatus int

const (
	StatusQueued FlagStatus = iota
	StatusSkipped
	StatusAccepted
	StatusRejected
)

func (s FlagStatus) String() string {
	switch s {
	case StatusAccepted:
		return "ACCEPTED"
	case StatusRejected:
		return "REJECTED"
	case StatusSkipped:
		return "SKIPPED"
	case StatusQueued:
		return "QUEUED"
	default:
		return "UNKNOWN"
	}
}

type SubmitResult struct {
	Flag   string
	Status FlagStatus
	Msg    string
}

var ResponsesMap = map[FlagStatus][]string{
	StatusAccepted: {"accepted", "congrat"},
	StatusRejected: {"bad", "wrong", "expired", "unknown", "your own",
		"too old", "not in database", "already submitted", "invalid flag"},
	StatusQueued: {"timeout", "game not started", "try again later", "game over", "is not up",
		"no such flag"},
}

type Team struct {
	Ip       string
	Name     string
}

type Flag struct {
    gorm.Model
    Flag     string     `gorm:"column:flag"`
    Sploit   string     `gorm:"column:sploit;index:flags_sploit"`
    Team     string     `gorm:"column:team;index:flags_team"`
    Time     time.Time  `gorm:"column:time;index:flags_time"`
    Status   FlagStatus `gorm:"column:status;index:flags_status_time,priority:1"`
    Response string     `gorm:"column:response"`
}

type Feeder struct {
	Enabled      bool
	Team        string
	Periodicity int
	PseudoChecker string
	Counter int
}
