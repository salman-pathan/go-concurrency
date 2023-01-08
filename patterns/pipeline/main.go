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
		"amazon.com",
		"reddit.com",
		"ebay.com",
		"google.com",
		"github.com",
		"stackoverflow.com",
		"gitlab.com",
		"msn.com",
		"duckduckgo.com",
		"yahoo.com",
		"yandex.ru",
	}

	doneChan := make(chan bool)
	defer close(doneChan)

	urlStream := urlStreamGenerator(doneChan, urls)
	pipeline := processUrls(doneChan, addPrefixGenerator(doneChan, urlStream))

	for resp := range pipeline {
		fmt.Printf("The url '%s' has returned '%d' status code\n", resp.Url, resp.StatusCode)
	}

}

func urlStreamGenerator(done <-chan bool, urls []string) <-chan string {
	urlStream := make(chan string)

	go func() {
		defer close(urlStream)

		for _, url := range urls {
			select {
			case <-done:
				return
			case urlStream <- url:
			}
		}
	}()

	return urlStream
}

func addPrefixGenerator(done <-chan bool, urlStream <-chan string) <-chan string {
	prefixStream := make(chan string)

	go func() {
		defer close(prefixStream)

		for url := range urlStream {
			prefixUrl := fmt.Sprintf("https://%s", url)
			select {
			case <-done:
				return
			case prefixStream <- prefixUrl:
			}
		}
	}()

	return prefixStream
}

func processUrls(done <-chan bool, urlStream <-chan string) <-chan Resp {
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
			case <-done:
				return
			case processStream <- r:
			}
		}
	}()

	return processStream
}
