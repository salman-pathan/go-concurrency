package main

import (
	"fmt"
	"sync"
)

type DBConnection struct{}

type DBManager struct {
	conn *DBConnection
}

func (dbMan *DBManager) getConn(wg *sync.WaitGroup) {
	defer wg.Done()
	if dbMan.conn == nil {
		fmt.Println("not initialized")
		dbMan.conn = &DBConnection{}
	}
	fmt.Println("initialized")
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
