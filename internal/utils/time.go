package utils

import "time"

type TimeUtil interface {
	Until(t time.Time) time.Duration
}

type RealTime struct{}

func (rt *RealTime) Until(t time.Time) time.Duration {
	return time.Until(t)
}
