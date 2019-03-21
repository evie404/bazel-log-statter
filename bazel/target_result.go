package bazel

import "time"

type Status string

const (
	StatusNoStatus Status = "NO STATUS"
	StatusPassed   Status = "PASSED"
	StatusUnknown  Status = "UNKNOWN"
	StatusFlaky    Status = "FLAKY"
)

type TargetResult struct {
	Name   string
	Cached bool
	Status
	time.Duration

	// flaky test attempts
	Successes int
	Attempts  int
}