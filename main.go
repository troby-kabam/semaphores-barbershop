package main

import (
	"fmt"
	"log"
	"sync"
	"time"
	"math/rand"
)

const (
	MAX_CUSTOMERS	= 12
	MAX_SEATS	= 5
)

type Barber struct {
	Cut		chan int
	Done		chan bool
}

func (b *Barber) Init() {
	b.Cut		= make(chan int, MAX_SEATS)
	b.Done		= make(chan bool)
}

func (b *Barber) Run() {
	for {
		x := <-b.Cut
		if x != 0 {
			fmt.Printf("Barber is cutting hair for customer %d.\n", x)
			b.CutHair()
		} else {
			b.Done <- true
			return
		}
	}
}

func (b *Barber) CutHair() {
	rand.Seed(time.Now().UnixNano())
	seed := rand.Intn(45) + 5
	duration, err := time.ParseDuration(fmt.Sprintf("%dms", seed * 10))
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(duration)
}

type Customer struct {
	Id		int
	Channel		chan int
	WaitGroup	*sync.WaitGroup
}

func (c *Customer) New(id int, ch chan int, wg *sync.WaitGroup) {
	c.Id = id
	c.Channel = ch
	c.WaitGroup = wg
}

func (c *Customer) GetHaircut() {
	defer c.WaitGroup.Done()
	select {
	case c.Channel <- c.Id:
		fmt.Printf("Customer %d has found a seat.\n", c.Id)
	case <-time.After(1 * time.Second):
		fmt.Printf("Customer %d balks!\n", c.Id)
	}
}

func main() {
	b := new(Barber)
	wg := new(sync.WaitGroup)

	b.Init()
	fmt.Println("Barbershop is now open.")
	go b.Run()
	wg.Add(MAX_CUSTOMERS)
	for i := 1; i <= MAX_CUSTOMERS; i++ {
		fmt.Printf("Customer %d has entered.\n", i)
		c := new(Customer)
		c.New(i, b.Cut, wg)
		go func() {
			c.GetHaircut()
		}()
	}

	wg.Wait()
	close(b.Cut)
	<-b.Done
	fmt.Println("The barbershop is now closed.")
}
