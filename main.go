package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

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

	for _, results := range allResults {
		for _, result := range results {
			fmt.Printf("%+v\n", result)
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
