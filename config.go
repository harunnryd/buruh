package buruh

import (
	"time"
)

type Config struct {
	MaxWorkerNum  uint64
	MinWorkerNum  uint64
	MaxWorkerLife time.Duration
	Debug         bool
}
