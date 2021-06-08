package util

import (
	"strconv"
	"time"
)

// Returns a formatted string of current date of format 01 June 2021
func GetCurrDate() string {
	curr := time.Now()
	return strconv.Itoa(curr.Day()) + " " + curr.Month().String() + " " + strconv.Itoa(curr.Year())
}
