package agents

import (
	"math/rand"
	"time"

	"github.com/projects/threaded-company-simulation/config"
)

type PersonType int

const (
	PT_CEO      PersonType = 1
	PT_WORKER   PersonType = 2
	PT_CUSTOMER PersonType = 3
	PT_MACHINE  PersonType = 4
)

func RandomSleepDuration(pt PersonType) time.Duration {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	coef := r.Float64()
	switch pt {
	case PT_CEO:
		return time.Duration(coef * config.AVERAGE_CEO_DELAY)
	case PT_WORKER:
		return time.Duration(coef * config.AVERAGE_WORKER_DELAY)
	case PT_CUSTOMER:
		return time.Duration(coef * config.AVERAGE_CUSTOMER_DELAY)
	case PT_MACHINE:
		return time.Duration(coef * config.AVERAGE_MACHINE_DELAY)
	default:
		return 0
	}
}
