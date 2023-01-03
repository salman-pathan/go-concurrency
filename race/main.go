package main

import (
	"fmt"
	"sync"
)

/*
*	A data race situation occurs when one operation attempts to read
*	a variable while another concurrenct operation is attempting to
*	write the same variable.
 */

var bankBalance = 2000

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
	bankBalance += amount
}
