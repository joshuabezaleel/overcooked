package main

import (
	"log"
	"sync"
	"time"
)

// Channel is queue while the table to put chopped vegetables is not so let's just assume that we will continuously fill
// the channel with chopped vegetables with the same amount.

type Salad int

const (
	CabbageSalad Salad = iota
	TomatoSalad
)

type Thing int

const (
	RawCabbage Thing = iota
	RawTomato
	ChoppedCabbage
	ChoppedTomato
)

type TableManager struct {
	Table    []Thing
	Capacity int
	mu       sync.Mutex
}

const TABLE_CAPACITY = 6

func main() {
	tableManager := &TableManager{Table: []Thing{}, Capacity: TABLE_CAPACITY}
	orders := []Salad{CabbageSalad, CabbageSalad, TomatoSalad, TomatoSalad}
	orderChan := make(chan Salad, len(orders))
	availablePut := make(chan bool, 1)

	for _, order := range orders {
		orderChan <- order
	}
	close(orderChan)

	for order := range orderChan {
		index := tableManager.isNotFull(availablePut)
		if index != -999 {
			var ingredient Thing
			switch order {
			case CabbageSalad:
				ingredient = RawCabbage
			case TomatoSalad:
				ingredient = RawTomato
			}
			tableManager.Put(index, ingredient, availablePut)
		}
	}
}

func chop(vegetable Thing) Thing {
	if vegetable == ChoppedCabbage || vegetable == ChoppedTomato {
		log.Printf("vegetable is chopped already!")
		return vegetable
	}

	log.Printf("chopping %v\n", vegetable)
	time.Sleep(time.Millisecond * 100)
	return vegetable
}

// isNotFull returns true when there is an empty slot to put a Thing on the table, and also
// the index of the table slice on where to put the thing. returns -999 and false otherwise.
func (tm *TableManager) isNotFull(availPut chan bool) int {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	for i, thing := range tm.Table {
		if thing == 0 {
			availPut <- true
			return i
		}
	}
	return -999
}

func (tm *TableManager) Put(index int, order Thing, availPut chan bool) {
	<-availPut
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.Table[index] = order
}

func (tm *TableManager) Get() error {
	return nil
}
