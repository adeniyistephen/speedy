package dlog

import (
	"time"
)

// Log is a single log object
type Log struct {
	Rank string `json:"rank"`

	Name   string `json:"name"`
	Caller string `json:"caller"` // function name

	Timestamp time.Time     `json:"timestamp"`
	Delta     time.Duration `json:"delta"`

	Data []interface{} `json:"data"`
}

// NewLog .
func NewLog(rank, name string, now time.Time, delta time.Duration, arg ...interface{}) *Log {
	return &Log{
		Rank: rank,

		Name:   name,
		Caller: getCaller(),

		Timestamp: now,
		Delta:     delta,

		Data: arg,
	}
}
