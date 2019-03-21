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
		var successes, failures, total int
		var totalDuration time.Duration

		for _, result := range results {
			switch result.Status {
			case bazel.StatusFailed:
				failures += 1
				total += 1
				totalDuration += result.Duration
			case bazel.StatusPassed:
				successes += 1
				total += 1
				totalDuration += result.Duration
			case bazel.StatusFlaky:
				successes += result.Successes
				failures += (result.Attempts - result.Successes)
				total += result.Attempts
				totalDuration += result.Duration
			}
		}

		if successes < total {
			successRatio := float64(successes) / float64(total)
			avgDuration := time.Duration(float64(totalDuration) / float64(total))

			fmt.Printf("%s: %v success %v\n", targetName, successRatio, avgDuration)
		}
	}
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
