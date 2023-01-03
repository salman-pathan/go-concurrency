package main

import "fmt"

var bankBalance = 1000

func main() {
	balance := make(chan int)

	go addBalance(balance)

	//	waiting to receive updated balance
	updatedBalance := <-balance
	fmt.Println("The balance is : ", updatedBalance)
}

func addBalance(balance chan int) {
	fmt.Println("The balance is : ", bankBalance)
	balance <- bankBalance + 100
	fmt.Println("Reachable statement")
}
