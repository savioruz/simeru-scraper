package utils

import "fmt"

func IntervalToCron(interval int) string {
	return fmt.Sprintf("0 */%d * * *", interval)
}
