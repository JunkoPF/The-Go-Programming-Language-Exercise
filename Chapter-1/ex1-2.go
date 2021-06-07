package main

import (
	"fmt"
	"os"
)

func main() {
	for index, value := range os.Args[1:] {
		fmt.Printf("Args[%d]: %s\n", index+1, value)
	}
}
