package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rickypai/bazel-log-statter/parser"
)

func main() {
	f, err := os.Open(buildFilePath(22141))
	defer f.Close()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		result, _ := parser.ParseLine(scanner.Text())

		if result != nil {
			fmt.Printf("%+v\n", result)
		}
	}

	println("lol")
}

func buildFilePath(id int) string {
	return filepath.Join("/Users/ricky/workspace/godel-logs", fmt.Sprintf("%v.txt", id))
}
