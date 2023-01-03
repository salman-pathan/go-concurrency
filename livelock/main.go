package main

import (
	"fmt"
	"sync"
	"time"
)

/*
*	Livelocks are programs that are actively performing
*	concurrent operations, but these operations do nothing
*	to move the state of the program forward.
 */

type Diner struct {
	Name     string
	IsHungry bool
}

func NewDiner(name string, isHungry bool) *Diner {
	return &Diner{Name: name, IsHungry: isHungry}
}

type Spoon struct {
	Owner *Diner
	mu    *sync.Mutex
}

func NewSpoon(owner *Diner) *Spoon {
	return &Spoon{Owner: owner}
}

func (s *Spoon) SetOwner(owner *Diner) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Owner = owner
}

func (s *Spoon) Use() {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Printf("%s has eaten!", s.Owner.Name)
}

func (d *Diner) EatWith(spoon *Spoon, spouse *Diner) {
	for {
		if !d.IsHungry {
			break
		}

		//	Don't have the spoon, wait until the spouse is done
		if spoon.Owner != d {
			time.Sleep(1 * time.Second)
			continue
		}

		//	If spouse is hungry, pass the spoon
		if spouse.IsHungry {
			fmt.Printf("%s : You eat first %s\n", d.Name, spouse.Name)
			spoon.Owner = spouse
			continue
		}

		//	When spouse is not hungry, you own the spoon
		spoon.Use()
		d.IsHungry = false

		//	Once done, pass the spoon to spouse
		spoon.Owner = spouse
	}
}

func main() {
	husband := NewDiner("John", true)
	wife := NewDiner("Alice", true)

	spoon := NewSpoon(husband)

	go husband.EatWith(spoon, wife)
	go wife.EatWith(spoon, husband)

	time.Sleep(10 * time.Second)
}
