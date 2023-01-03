package main

import (
	"fmt"
	"sync"
)

/*
*	Solving the issue with Memory Access Syncronization
 */

var bankBalance = 2000
var mu sync.Mutex

func main() {
	wg := sync.WaitGroup{}
	for i := 1; i < 31; i++ {
		wg.Add(1)
		go addBalance(i, &wg)
	}

	wg.Wait()
	fmt.Println("Bank Balance : ", bankBalance)
}

func addBalance(amount int, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	defer mu.Unlock()
	bankBalance += amount
}
