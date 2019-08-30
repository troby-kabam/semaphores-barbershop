package main

import (
	"fmt"
	"time"
)

type Barber struct {
	Cut	chan int
}

func (b *Barber) Init() {
	b.Cut = make(chan int)
}

func (b *Barber) Run() {
	for {
		fmt.Println("Barber is sleeping.")
		x := <- b.Cut
		fmt.Printf("Barber is cutting hair for customer %d.\n", x)
		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {
	b := new(Barber)
	b.Init()
	go b.Run()
	b.Cut <- 1
	b.Cut <- 2
	b.Cut <- 3
	

/*
	customer enters
		if barber asleep
			wake and getHaircut
		if barber busy
			try chair
		if chair full
			leave
*/
}
