package utils

import "time"

type TimeUtil interface {
	Until(t time.Time) time.Duration
}
