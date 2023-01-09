package main

import (
	"fmt"
	"net/http"
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
	worker1 := processUrlStream(doneChan, urlStream)
	worker2 := processUrlStream(doneChan, urlStream)
	worker3 := processUrlStream(doneChan, urlStream)

	for resp := range worker1 {
		fmt.Printf("The url '%s' has returned '%d' status code. Worker 1\n", resp.Url, resp.StatusCode)
	}

	for resp := range worker2 {
		fmt.Printf("The url '%s' has returned '%d' status code. Worker 2\n", resp.Url, resp.StatusCode)
	}

	for resp := range worker3 {
		fmt.Printf("The url '%s' has returned '%d' status code. Worker 3\n", resp.Url, resp.StatusCode)
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

func processUrlStream(doneChan <-chan bool, urlStream <-chan string) <-chan Resp {
	processStream := make(chan Resp)

	go func() {
		defer close(processStream)
		for url := range urlStream {
			resp, err := http.Get(url)
			if err != nil {
				return
			}
			r := Resp{Url: url, StatusCode: resp.StatusCode}
			select {
			case <-doneChan:
				return
			case processStream <- r:
			}
		}
	}()

	return processStream
}
