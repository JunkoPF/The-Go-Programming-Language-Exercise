package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	files := os.Args[1:]
	count := make(map[string]int)

	if len(files) == 0 {
		countLines(os.Stdin, &count)
	} else {
		for _, file := range files {
			fd, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup: %v\n", err)
				continue
			}
			countLines(fd, &count)
		}
	}

	flag := true

	for line, n := range count {
		if n > 1 {
			if flag {
				fmt.Printf("Duplications exists in such files: %v\n", strings.Join(files, ", "))
				flag = false
			}
			fmt.Printf("%v\t%v\n", line, n)
		}
	}
}

func countLines(fd *os.File, count *map[string]int) {
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		(*count)[scanner.Text()]++
	}
}
