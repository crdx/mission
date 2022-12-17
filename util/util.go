package util

import (
	"time"
)

func Pluralise(n int, singular, plural string) string {
	if n == 1 {
		return singular
	} else {
		return plural
	}
}

func GetRunType(headless bool) string {
	if headless {
		return "auto"
	} else {
		return "manual"
	}
}

func FormatDuration(duration time.Duration) string {
	return duration.Truncate(time.Millisecond).String()
}

func Contains[T comparable](items []T, item T) bool {
	for _, x := range items {
		if x == item {
			return true
		}
	}
	return false
}
