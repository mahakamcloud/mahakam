package util

import (
	"math/rand"
	"time"
)

func RandomIntRange(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn((max-min)+1) + min
}
