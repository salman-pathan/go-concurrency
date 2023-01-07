package main

import (
	"fmt"
	"net/http"
	"sync"
)

type Resp struct {
	Url        string
	StatusCode int
}

func main() {

	urls := []string{
		"https://amazon.com",
		"https://reddit.com",
		"https://ebay.com",
		"https://google.com",
		"https://github.com",
		"https://stacckoverflow.com",
		"http://gitlab.com",
		"https://msn.com",
		"https://duckduckgogo.com",
		"https://yahoo.com",
		"https://yandex.ru",
	}

	doneChan := make(chan bool)
	respChan := make(chan Resp)
	wg := sync.WaitGroup{}

	for _, url := range urls {
		wg.Add(1)
		go processUrl(url, respChan, doneChan, &wg)
	}

	go func() {
		wg.Wait()
		close(respChan)
	}()

loop:
	for {
		select {
		case resp, ok := <-respChan:
			if !ok {
				return
			}
			fmt.Printf("The url '%s' return '%d' status code\n", resp.Url, resp.StatusCode)
		case done, ok := <-doneChan:
			if !ok || done {
				break loop
			}
		}
	}

}

func processUrl(url string, respChan chan<- Resp, done chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		done <- true
		return
	}
	respChan <- Resp{
		Url:        url,
		StatusCode: resp.StatusCode,
	}
}
