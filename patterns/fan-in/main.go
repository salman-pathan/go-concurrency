package main

import (
	"fmt"
	"net/http"
	"sync"
)

type Resp struct {
	Url        string
	StatusCode int
	WorkerId   int
}

func main() {

	urls := []string{
		"https://amazon.com",
		"https://reddit.com",
		"https://ebay.com",
		"https://google.com",
		"https://github.com",
		"https://stackoverflow.com",
		"https://gitlab.com",
		"https://msn.com",
		"https://duckduckgo.com",
		"https://yahoo.com",
		"https://yandex.ru",
	}

	doneChan := make(chan bool)
	defer close(doneChan)

	urlStream := generateUrlStream(doneChan, urls)

	//	Fanning out
	worker1 := processUrlStream(doneChan, urlStream, 1)
	worker2 := processUrlStream(doneChan, urlStream, 2)
	worker3 := processUrlStream(doneChan, urlStream, 3)

	//	Fanning in
	respStream := fanIn(doneChan, worker1, worker2, worker3)

	for resp := range respStream {
		fmt.Printf("The url '%s' has returned '%d' status code. Worker %d\n", resp.Url, resp.StatusCode, resp.WorkerId)
	}

}

func generateUrlStream(doneChan <-chan bool, urls []string) <-chan string {
	urlStream := make(chan string)

	go func() {
		defer close(urlStream)

		for _, url := range urls {
			select {
			case <-doneChan:
				return
			case urlStream <- url:
			}
		}
	}()

	return urlStream
}

func processUrlStream(doneChan <-chan bool, urlStream <-chan string, workerID int) <-chan Resp {
	processStream := make(chan Resp)

	go func() {
		defer close(processStream)
		for url := range urlStream {
			resp, err := http.Get(url)
			if err != nil {
				return
			}
			r := Resp{Url: url, StatusCode: resp.StatusCode, WorkerId: workerID}
			select {
			case <-doneChan:
				return
			case processStream <- r:
			}
		}
	}()

	return processStream
}

func fanIn(doneChan <-chan bool, urlChans ...<-chan Resp) <-chan Resp {
	var wg sync.WaitGroup
	respStream := make(chan Resp)

	wg.Add(len(urlChans))
	for _, urlChan := range urlChans {

		go func(urlChan <-chan Resp) {
			defer wg.Done()
			for url := range urlChan {
				select {
				case <-doneChan:
					return
				case respStream <- url:
				}
			}
		}(urlChan)

	}

	go func() {
		wg.Wait()
		defer close(respStream)
	}()

	return respStream
}
