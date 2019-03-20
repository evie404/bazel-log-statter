package parser

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Status string

var (
	cachedLineRegex   = regexp.MustCompile(`(?P<target>\/\/.+)\s+(?P<cached>\(cached\))\s+(?P<status>PASSED)\s+in\s+(?P<duration>.+)s`)
	uncachedLineRegex = regexp.MustCompile(`(?P<target>\/\/.+)\s+(?P<status>PASSED|FAILED)\s+in\s+(?P<duration>.+)s`)
	noStatusLineRegex = regexp.MustCompile(`(?P<target>\/\/.+)\s+(?P<status>NO\sSTATUS)`)
	flakyLineRegex    = regexp.MustCompile(`(?P<target>\/\/.+)\s+(?P<status>FLAKY),\sfailed\sin(?P<success>.+)\sout\sof(?P<tries>.+)in\s+(?P<duration>.+)s`)
)

const (
	StatusNoStatus Status = "NO STATUS"
	StatusPassed   Status = "PASSED"
	StatusUnknown  Status = "UNKNOWN"
	StatusFlaky    Status = "FLAKY"
)

func ParseLine(line string) (target string, cached bool, status Status, duration time.Duration, err error) {
	var matched bool

	matched, target, cached, status, duration, err = cachedMatches(line)
	if matched || err != nil {
		return
	}

	matched, target, cached, status, duration, err = uncachedMatches(line)
	if matched || err != nil {
		return
	}

	matched, target, cached, status, duration, err = noStatusMatches(line)
	if matched || err != nil {
		return
	}

	return "", false, StatusUnknown, 0, nil
}

func cachedMatches(line string) (matched bool, target string, cached bool, status Status, duration time.Duration, err error) {
	matches := cachedLineRegex.FindStringSubmatch(line)

	if len(matches) == 0 {
		return false, "", false, StatusUnknown, 0, err
	}

	target = strings.TrimSpace(matches[1])
	cached = matches[2] == "(cached)"
	status = Status(matches[3])
	duration, err = parseDuration(matches[4])
	if err != nil {
		return false, "", false, StatusUnknown, 0, err
	}

	return true, target, cached, status, duration, nil
}

func uncachedMatches(line string) (matched bool, target string, cached bool, status Status, duration time.Duration, err error) {
	matches := uncachedLineRegex.FindStringSubmatch(line)

	if len(matches) == 0 {
		return false, "", false, StatusUnknown, 0, err
	}

	target = strings.TrimSpace(matches[1])
	status = Status(matches[2])
	duration, err = parseDuration(matches[3])
	if err != nil {
		return false, "", false, StatusUnknown, 0, err
	}

	return true, target, false, status, duration, nil
}

func noStatusMatches(line string) (matched bool, target string, cached bool, status Status, duration time.Duration, err error) {
	matches := noStatusLineRegex.FindStringSubmatch(line)

	if len(matches) == 0 {
		return false, "", false, StatusUnknown, 0, err
	}

	target = strings.TrimSpace(matches[1])
	status = Status(matches[2])

	return true, target, false, status, 0, nil
}

func parseDuration(durationStr string) (time.Duration, error) {
	durationF, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, err
	}
	return time.Duration(durationF * float64(time.Second)), nil
}
