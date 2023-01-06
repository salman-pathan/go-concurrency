package main

import (
	"fmt"
	"sync"
)

type DBConnection struct{}

type DBManager struct {
	conn *DBConnection
	once sync.Once
}

func (dbMan *DBManager) getConn(wg *sync.WaitGroup) {
	defer wg.Done()
	dbMan.once.Do(func() {
		fmt.Println("not initialized")
		dbMan.conn = &DBConnection{}
	})
	if dbMan.conn != nil {
		fmt.Println("initialized")
	}
}

func main() {
	dbMan := DBManager{}
	wg := sync.WaitGroup{}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go dbMan.getConn(&wg)
	}
	wg.Wait()
}
