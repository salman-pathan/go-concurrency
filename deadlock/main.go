package main

import "fmt"

/*
*	In a deadlock situation all concurrent processes
*	wait on one another, and the program will never
*	recover without outside intervention.
 */

var bankBalance = 1000

func main() {
	balance := make(chan int)

	go addBalance(balance)
	//	waiting to receive updated balance
	<-balance
}

func addBalance(balance chan int) {
	fmt.Println("The balance is : ", bankBalance)
	<-balance
	fmt.Println("Unreachable statement")
}
