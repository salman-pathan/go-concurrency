package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

var inventory = make(map[string]int)

func main() {
	mu := sync.Mutex{}
	cond := sync.NewCond(&mu)

	wg := sync.WaitGroup{}
	wg.Add(4)

	go buy(cond, &wg, "John", "Apple", 2)
	go buy(cond, &wg, "Alice", "Banana", 4)
	go buy(cond, &wg, "Jenny", "Pineapple", 1)
	go buy(cond, &wg, "John", "Banana", 4)

	wg.Wait()

	stockUp(cond)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

func stockUp(cond *sync.Cond) {
	cond.L.Lock()
	defer cond.L.Unlock()

	fmt.Println("adding apples")
	inventory["Apple"] = 10
	time.Sleep(1 * time.Second)

	fmt.Println("adding bananas")
	inventory["Banana"] = 20
	time.Sleep(1 * time.Second)

	fmt.Println("adding pineapple")
	inventory["Pineapple"] = 5
	time.Sleep(1 * time.Second)

	fmt.Println("broadcasting")
	cond.Broadcast()
}

func buy(cond *sync.Cond, wg *sync.WaitGroup, customerName string, fruit string, quantity int) {
	cond.L.Lock()
	wg.Done()
	fmt.Printf("Customer %s is waiting to buy %d %s\n", customerName, quantity, fruit)
	defer cond.L.Unlock()
	cond.Wait()
	inventory[fruit] -= quantity
	fmt.Printf("Customer %s has bought %d %s\n", customerName, quantity, fruit)
}
