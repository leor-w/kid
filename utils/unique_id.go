package utils

import (
	"math/rand"
	"time"
)

// UniqueId 获取新的唯一 ID
func UniqueId(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	uniqueId := rand.Int63n(max)
	if uniqueId < min {
		return min + uniqueId
	}
	return uniqueId
}
