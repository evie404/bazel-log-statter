package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rickypai/bazel-log-statter/bazel"
	"github.com/rickypai/bazel-log-statter/parser"
)

func main() {
	var startBuild, endBuild int
	var sortMethod string
	var ignoreCached bool

	flag.IntVar(&startBuild, "start", 0, "start")
	flag.IntVar(&endBuild, "end", 0, "end")
	flag.StringVar(&sortMethod, "sort", "name", "sort")
	flag.BoolVar(&ignoreCached, "ignore-cached", false, "ignore-cached")

	flag.Parse()

	if startBuild == 0 {
		panic("-start flag required")
	}

	if endBuild == 0 {
		panic("-end flag required")
	}

	builds := endBuild - startBuild + 1

	allResults := make([][]*bazel.TargetResult, builds+1)

	var wg sync.WaitGroup

	for i := startBuild; i <= endBuild; i++ {
		wg.Add(1)

		go func(fileNum, index int) {
			defer wg.Done()

			var parseErr error
			// println(fileNum)
			allResults[index], parseErr = parseFile(fileNum)
			if parseErr != nil {
				println(parseErr)
			}
		}(i, i-startBuild)
	}

	wg.Wait()

	targetResults := map[string]*AggregateResult{}

	for _, results := range allResults {
		for _, result := range results {
			if ignoreCached && result.Cached {
				continue
			}

			if _, found := targetResults[result.Name]; !found {
				targetResults[result.Name] = &AggregateResult{
					TargetName: result.Name,
				}
			}
			targetResults[result.Name].AddResult(result)
		}
	}

	targetNames := make([]string, 0, len(targetResults))
	longestNameLen := 0
	longestTriesLen := 0

	for targetName, aggregate := range targetResults {
		targetNames = append(targetNames, targetName)
		if !aggregate.AllSuccesses() {
			if len(targetName) > longestNameLen {
				longestNameLen = len(targetName)
			}

			if len(strconv.Itoa(aggregate.Total)) > longestTriesLen {
				longestTriesLen = len(strconv.Itoa(aggregate.Total))
			}
		}
	}

	switch sortMethod {
	case "name":
		sort.Strings(targetNames)
	case "failures":
		sort.Slice(targetNames, func(i, j int) bool {
			return targetResults[targetNames[i]].SuccessRatio() < targetResults[targetNames[j]].SuccessRatio()
		})
	case "successes":
		sort.Slice(targetNames, func(i, j int) bool {
			return targetResults[targetNames[i]].SuccessRatio() > targetResults[targetNames[j]].SuccessRatio()
		})
	case "longest":
		sort.Slice(targetNames, func(i, j int) bool {
			return targetResults[targetNames[i]].AverageDuration() > targetResults[targetNames[j]].AverageDuration()
		})
	case "shortest":
		sort.Slice(targetNames, func(i, j int) bool {
			return targetResults[targetNames[i]].AverageDuration() < targetResults[targetNames[j]].AverageDuration()
		})
	default:
		sort.Strings(targetNames)
	}

	for _, targetName := range targetNames {
		aggregate := targetResults[targetName]

		if !aggregate.AllSuccesses() {
			spaces := strings.Join(make([]string, 2+longestNameLen-len(aggregate.TargetName)), " ")
			triesSpaces := strings.Join(make([]string, 1+longestTriesLen-len(strconv.Itoa(aggregate.Total))), " ")

			fmt.Printf("%s%.2f%% success in %v%v tries %v\n", aggregate.TargetName+spaces, aggregate.SuccessRatio(), triesSpaces, aggregate.Total, aggregate.AverageDuration())
		}
	}
}

type AggregateResult struct {
	TargetName    string
	Targets       []*bazel.TargetResult
	Total         int
	Successes     int
	Failures      int
	TotalDuration time.Duration
}

func (a *AggregateResult) AddResult(result *bazel.TargetResult) {
	switch result.Status {
	case bazel.StatusFailed:
		a.Failures += 1
		a.Total += 1
		a.TotalDuration += result.Duration
	case bazel.StatusPassed:
		a.Successes += 1
		a.Total += 1
		a.TotalDuration += result.Duration
	case bazel.StatusFlaky:
		a.Successes += result.Successes
		a.Failures += (result.Attempts - result.Successes)
		a.Total += result.Attempts
		a.TotalDuration += result.Duration
	}
}

func (a *AggregateResult) AllSuccesses() bool {
	return a.Successes == a.Total
}

func (a *AggregateResult) SuccessRatio() float64 {
	return float64(a.Successes*100) / float64(a.Total)
}

func (a *AggregateResult) AverageDuration() time.Duration {
	return time.Duration(float64(a.TotalDuration) / float64(a.Total))
}

func parseFile(id int) ([]*bazel.TargetResult, error) {
	f, err := os.Open(buildFilePath(id))
	defer f.Close()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)

	results := make([]*bazel.TargetResult, 0)

	for scanner.Scan() {
		result, _ := parser.ParseLine(scanner.Text())

		if result != nil {
			results = append(results, result)
		}
	}

	return results, nil
}

func buildFilePath(id int) string {
	return filepath.Join("/Users/ricky/workspace/godel-logs", fmt.Sprintf("%v.txt", id))
}
