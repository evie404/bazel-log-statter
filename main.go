package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rickypai/bazel-log-statter/bazel"
	"github.com/rickypai/bazel-log-statter/parser"
)

func main() {
	startBuild := 21000
	endBuild := 22227
	builds := endBuild - startBuild + 1

	allResults := make([][]*bazel.TargetResult, builds+1)

	for i := startBuild; i <= endBuild; i++ {
		// go func(i int) {
		var parseErr error
		println(i)
		allResults[i-startBuild], parseErr = parseFile(i)
		if parseErr != nil {
			println(parseErr)
		}
		// }(i)
	}

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
