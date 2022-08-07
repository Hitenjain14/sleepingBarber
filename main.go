package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var seatingCapacity = 10
var arrivalRate = 100
var cutDuration = 1000 * time.Millisecond
var timeOpen = 10 * time.Second

func main() {

	// seed our random number generator

	rand.Seed(time.Now().UnixNano())

	color.Yellow("Sleeping Barber problem")
	color.Yellow("--------------------------------")

	//create channels
	clienChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	//create barbershop

	shop := BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		BarbersDoneChan: doneChan,
		ClientsChan:     clienChan,
		Open:            true,
	}

	color.Green("Shop is open for the day!")

	// add barbers
	shop.addBarber("Frank")
	shop.addBarber("Tom")
	shop.addBarber("Harry")
	shop.addBarber("Joe")
	shop.addBarber("Bill")

	// start the barbershop as a goroutine

	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(timeOpen)
		shopClosing <- true
		shop.closeShopForTheDay()
		closed <- true
	}()

	// add clients
	i := 1
	go func() {
		for {
			randomMillisecond := rand.Int() % 2 * arrivalRate
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMillisecond)):
				shop.addClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()

	//block until barbershop is closed

	<-closed

}
