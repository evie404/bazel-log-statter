package parser

import "time"

type Status string

const (
	Status_NoStatus = "NO STATUS"
	Status_Passed   = "PASSED"
)

func ParseLine(line string) (target string, cached bool, status string, duration time.Duration) {
	return "", false, "", time.Duration(1)
}
