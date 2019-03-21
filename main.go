package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rickypai/bazel-log-statter/bazel"
	"github.com/rickypai/bazel-log-statter/parser"
)

func main() {
	var startBuild, endBuild int

	flag.IntVar(&startBuild, "start", 0, "start")
	flag.IntVar(&endBuild, "end", 0, "end")
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

	targetResults := map[string][]*bazel.TargetResult{}

	for _, results := range allResults {
		for _, result := range results {
			if _, found := targetResults[result.Name]; !found {
				targetResults[result.Name] = []*bazel.TargetResult{}
			}

			targetResults[result.Name] = append(targetResults[result.Name], result)
		}
	}

	for targetName, results := range targetResults {
		aggregate := &AggregateResult{
			TargetName: targetName,
		}

		for _, result := range results {
			switch result.Status {
			case bazel.StatusFailed:
				aggregate.Failures += 1
				aggregate.Total += 1
				aggregate.TotalDuration += result.Duration
			case bazel.StatusPassed:
				aggregate.Successes += 1
				aggregate.Total += 1
				aggregate.TotalDuration += result.Duration
			case bazel.StatusFlaky:
				aggregate.Successes += result.Successes
				aggregate.Failures += (result.Attempts - result.Successes)
				aggregate.Total += result.Attempts
				aggregate.TotalDuration += result.Duration
			}
		}

		if !aggregate.AllSuccesses() {
			fmt.Printf("%s: %.2f%% success %v\n", aggregate.TargetName, aggregate.SuccessRatio(), aggregate.AverageDuration())
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
