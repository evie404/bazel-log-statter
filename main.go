package main

import (
	"bufio"
	"os"
	"path/filepath"
)

func main() {
	f, err := os.Open(filepath.Join("/Users/ricky/workspace/godel-logs", "22141.txt"))
	defer f.Close()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		println(scanner.Text())
	}

	println("lol")
}
