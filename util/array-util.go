package util

import "time"

func Includes(arr []int, val int) bool {
	for _, e := range arr {
		if e == val {
			return true
		}
	}
	return false
}

func GetNextLocaleDay(locale string) string {
	loc, _ := time.LoadLocation(locale)
	now := time.Now().In(loc)
	now = now.AddDate(0, 0, 1)
	tomorrowDateIST := now.Format("02-01-2006")
	return tomorrowDateIST
}
