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

	for resp := range orDone(doneChan, respChan) {
		fmt.Printf("The url '%s' has returned '%d' status code\n", resp.Url, resp.StatusCode)
	}

}

func processUrl(url string, respChan chan<- Resp, done chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		done <- true
		return
	}
	respChan <- Resp{
		Url:        url,
		StatusCode: resp.StatusCode,
	}
}

func orDone(done <-chan bool, c <-chan Resp) <-chan Resp {
	valStream := make(chan Resp)
	go func() {
		defer close(valStream)

		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}
