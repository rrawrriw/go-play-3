package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	maxComputers = 4
)

var (
	heros = []string{
		"Alice",
		"Bob",
		"Charlie",
		"Robocop",
		"Batman",
		"Spiderman",
		"Superman",
		"Antman",
		"Auqaman",
		"Hulk",
		"Ironman",
	}

	wg sync.WaitGroup
)

type (
	cafe struct {
		name  string
		next  []chan struct{}
		queue chan struct{}
	}
)

func main() {
	coffeeShop := newCafe("crazy horse")

	for _, hero := range heros {
		coffeeShop.add(hero)
	}

	coffeeShop.open()
}

func newCafe(n string) cafe {
	return cafe{
		name:  n,
		next:  []chan struct{}{},
		queue: make(chan struct{}, maxComputers),
	}
}

func (c cafe) open() {
	fmt.Println("-- OPEN SHOP", c.name, "--")
	for _, n := range c.next {
		n <- struct{}{}
		c.queue <- struct{}{}
	}

	wg.Wait()
}

func (c *cafe) add(hero string) {
	cs := *c
	// Die Nummer entspricht dem Index im slice
	n := make(chan struct{})
	cs.next = append(cs.next, n)
	*c = cs

	wg.Add(1)
	go func(name string, queue chan struct{}, next chan struct{}) {
		defer wg.Done()
		fmt.Println("wake up", name)
		<-next
		fmt.Println("\t", name, "is online")
		sleep := d(700, 900)
		time.Sleep(sleep)
		fmt.Println("\t\t", name, "is done, was online for", sleep)
		<-queue
	}(hero, cs.queue, n)

}

func d(s, e int) time.Duration {
	now := time.Now().Nanosecond()
	src := rand.NewSource(int64(now))
	r := rand.New(src)
	n := r.Intn(e-s) + s

	return time.Duration(int64(n)) * time.Millisecond
}
