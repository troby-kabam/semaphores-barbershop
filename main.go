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
	BARBER_MAX	= 50    // (BARBER_MAX * 10) millisecond sleep time
	CUSTOMER_WAIT	= 1	// wait time in seconds to find a seat
)

type Barber struct {
	Cut		chan int	// customer id intake - buffered
	Done		chan bool	// telemetry
}

func NewBarber() *Barber {
	b := new(Barber)
	b.Cut		= make(chan int, MAX_SEATS)
	b.Done		= make(chan bool)
	return b
}

func (b *Barber) Run() {
	for {
		// barber sleeps here until a customer sends their id
		x, more := <-b.Cut
		// if input is 0 we do not have a valid customer id
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
	seed := rand.Intn(BARBER_MAX) + 5	// random value between 5 and BARBER_MAX
	if seed > BARBER_MAX {
		seed = BARBER_MAX
	}
	duration, err := time.ParseDuration(fmt.Sprintf("%dms", seed * 10)) // millisecond sleep time
	if err != nil {
		// no known reason to encounter this
		log.Fatal(err)
	}
	time.Sleep(duration)	// 5ms >= duration < (BARBER_MAX * 10) ms
}

type Customer struct {
	Id		int		// id number
	Channel		chan int	// Barber.Cut buffered MAX_SEATS
	WaitGroup	*sync.WaitGroup	// sync group
}

func NewCustomer(id int, ch chan int, wg *sync.WaitGroup) *Customer {
	return &Customer{id, ch, wg}
}

/*
 * Customer.GetHaircut uses a timeout of CUSTOMER_WAIT seconds to
 * determine whether a customer gets a spot in the barber's queue.
 * Customers who encounter the timeout will balk.
 */
func (c *Customer) GetHaircut() {
	defer c.WaitGroup.Done()
	select {
	case c.Channel <- c.Id:
		fmt.Printf("Customer %d has found a seat.\n", c.Id)
	case <-time.After(CUSTOMER_WAIT * time.Second):
		fmt.Printf("Customer %d balks!\n", c.Id)
	}
}

func main() {
	b := NewBarber()
	wg := new(sync.WaitGroup)

	fmt.Println("Barbershop is now open.")
	go b.Run()
	wg.Add(MAX_CUSTOMERS)
	for i := 1; i <= MAX_CUSTOMERS; i++ {
		fmt.Printf("Customer %d has entered.\n", i)
		c := NewCustomer(i, b.Cut, wg)
		go func() {
			c.GetHaircut()
		}()
	}

	wg.Wait()	// sync Customer threads
	close(b.Cut)	// close buffered channel
	<-b.Done	// sync Barber thread
	fmt.Println("The barbershop is now closed.")
}
