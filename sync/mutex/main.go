package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
*	Mutex stands for "mutual exclusion" and is a way to guard critical sections.
*	A mutex provides a concurrent-safe way to express exclusive access to shared
*	resources.
 */
func main() {
	fruitMap := map[string]int{
		"Apple":      10,
		"Banana":     24,
		"Mango":      4,
		"Pineapple":  14,
		"Orange":     16,
		"Watermelon": 3,
		"Papaya":     2,
		"Raspberry":  20,
		"Jackfruit":  1,
		"Nectarine":  6,
		"Cantaloupe": 2,
		"Lychee":     30,
	}

	wg := sync.WaitGroup{}
	wg.Add(3)

	go stockFruits(fruitMap, &wg)
	go sellFruits(fruitMap, &wg)
	go sellFruits(fruitMap, &wg)

	wg.Wait()
}

func stockFruits(fruits map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	min := 1
	max := 5
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 30; i++ {
		for k := range fruits {
			incr := rand.Intn((max - min) + min)
			fruits[k] = fruits[k] + incr
		}
	}
}

func sellFruits(fruits map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 30; i++ {
		for k, v := range fruits {
			if v != 0 {
				fruits[k] = fruits[k] - 1
			}
		}
	}
}

func fruitStats(fruits map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	for k, v := range fruits {
		fmt.Printf("%s current count is %d\n", k, v)
	}
}
