package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ShouldFail() bool {
	return rand.Float32() < 0.1
}

func ShouldDelay() bool {
	return rand.Float32() < 0.2
}

func RandomDelay() time.Duration {
	return time.Duration(500 + rand.Intn(1500)) * time.Millisecond
}
