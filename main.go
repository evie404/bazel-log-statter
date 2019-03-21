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
	results := parseFile(22141)

	for _, result := range results {
		fmt.Printf("%+v\n", result)
	}

	println("done")
}

func parseFile(id int) []*bazel.TargetResult {
	f, err := os.Open(buildFilePath(id))
	defer f.Close()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	results := make([]*bazel.TargetResult, 0)

	for scanner.Scan() {
		result, _ := parser.ParseLine(scanner.Text())

		if result != nil {
			results = append(results, result)
		}
	}

	return results
}

func buildFilePath(id int) string {
	return filepath.Join("/Users/ricky/workspace/godel-logs", fmt.Sprintf("%v.txt", id))
}
