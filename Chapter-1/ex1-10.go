package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	urls := os.Args[1:]
	ch := make(chan string)
	if len(urls) == 0 {
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			wg.Add(1)
			go func() {
				defer wg.Done()
				fetch(input.Text(), ch)
			}()
			wg.Add(1)
			go func(ch <-chan string) {
				defer wg.Done()
				fmt.Println(<-ch)
			}(ch)
		}
	} else {
		for _, url := range urls {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				fetch(url, ch)
			}(url)
		}
		for i := 0; i < len(urls); i++ {
			fmt.Println(<-ch)
		}
	}
	wg.Wait()

}

func fetch(url string, ch chan<- string) {
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("fetch: %v", err)
		return
	}
	output := bufio.NewScanner(resp.Body)
	cntlen := 0
	for output.Scan() {
		cntlen += len(output.Text())
	}

	cost := time.Since(start).Seconds()
	ch <- fmt.Sprintf("url = %v | total length = %v bytes | cost = %v sec", url, cntlen, cost)
}
