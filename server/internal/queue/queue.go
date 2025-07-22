package queue

import (
	"github.com/vert3xc/bebrochka/server/internal/models"
)

var FlagChan = make(chan models.Flag, 1000)
