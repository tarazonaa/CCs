package utils

import "time"

func GetCurrentTS() time.Time {
	return time.Now().UTC()
}