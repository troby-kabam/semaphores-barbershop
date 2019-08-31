package main

import (
	"fmt"
	"log"
	"sync"
	"time"
	"math/rand"
)

const (
	MAX_CUSTOMERS	= 12	// one thread per customer
	MAX_SEATS	= 6	// 5 waiting chairs + 1 barber chair
)

type Barber struct {
	Cut		chan int	// customer id intake - buffered
	Done		chan bool	// telemetry
}

func (b *Barber) Init() {
	b.Cut		= make(chan int, MAX_SEATS)
	b.Done		= make(chan bool)
}

func (b *Barber) Run() {
	for {
		// barber sleeps here until a customer sends their id
		x, more := <-b.Cut
		// if input is 0 we are done processing customers
		if x != 0 {
			fmt.Printf("Barber is cutting hair for customer %d.\n", x)
			b.CutHair()
		}
		if more == false {
			b.Done <- true
			return
		}
	}
}

func (b *Barber) CutHair() {
	rand.Seed(time.Now().UnixNano())
	seed := rand.Intn(45) + 5	// random value between 5 and 50
	duration, err := time.ParseDuration(fmt.Sprintf("%dms", seed * 10)) // 5 - 500 ms sleep time
	if err != nil {
		// no known reason to encounter this
		log.Fatal(err)
	}
	time.Sleep(duration)
}

type Customer struct {
	Id		int		// id number
	Channel		chan int	// Barber.Cut buffered MAX_SEATS
	WaitGroup	*sync.WaitGroup	// sync group
}

func (c *Customer) New(id int, ch chan int, wg *sync.WaitGroup) {
	c.Id = id
	c.Channel = ch
	c.WaitGroup = wg
}

/*
 * Customer.GetHaircut uses a timeout of 1 second
 * to determine whether a customer gets a spot
 * in the barber's queue. Customers who encounter
 * the timeout will balk.
 */
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

	wg.Wait()	// sync Customer threads
	close(b.Cut)	// close buffered channel
	<-b.Done	// sync Barber thread
	fmt.Println("The barbershop is now closed.")
}
