package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

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

	mut := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(2)

	go stockFruits(fruitMap, &wg, &mut)
	go sellFruits(fruitMap, &wg, &mut)

	wg.Wait()

	fruitStats(fruitMap)

}

func stockFruits(fruits map[string]int, wg *sync.WaitGroup, mut *sync.Mutex) {
	defer wg.Done()
	min := 1
	max := 5
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 30; i++ {
		for k := range fruits {
			mut.Lock()
			incr := rand.Intn((max - min) + min)
			fruits[k] = fruits[k] + incr
			mut.Unlock()
		}
	}
}

func sellFruits(fruits map[string]int, wg *sync.WaitGroup, mut *sync.Mutex) {
	defer wg.Done()
	for i := 1; i <= 30; i++ {
		for k, v := range fruits {
			if v != 0 {
				mut.Lock()
				fruits[k] = fruits[k] - 1
				mut.Unlock()
			}
		}
	}
}

func fruitStats(fruits map[string]int) {
	for k, v := range fruits {
		fmt.Printf("%s current count is %d\n", k, v)
	}
}
