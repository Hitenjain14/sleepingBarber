package main

import (
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientsChan     chan string
	Open            bool
}

func (shop *BarberShop) addBarber(barber string) {
	shop.NumberOfBarbers++

	go func() {
		isSleeping := false
		color.Yellow("%s goes to waiting room to check clients", barber)

		for {
			if len(shop.ClientsChan) == 0 {
				color.Yellow("There is nothing to do so %s is sleeping", barber)
				isSleeping = true
			}

			client, shopOpen := <-shop.ClientsChan

			if shopOpen {
				if isSleeping {
					color.Yellow("%s wakes %s up", barber, client)
					isSleeping = false
				}
				//cut hair
				shop.cutHair(client, barber)
			} else {
				//shop is closed so send the barber home
				shop.sendBarberHome(barber)
				return
			}

		}
	}()
}

func (shop *BarberShop) cutHair(client string, barber string) {
	color.Green("%s cuts hair for %s", barber, client)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s is done cutting hair for %s", barber, client)
}

func (shop *BarberShop) sendBarberHome(barber string) {
	color.Cyan("%s is going home", barber)
	shop.BarbersDoneChan <- true
}

func (shop *BarberShop) closeShopForTheDay() {
	color.Cyan("Closing shop for the day")
	close(shop.ClientsChan)
	shop.Open = false

	for a := 1; a <= shop.NumberOfBarbers; a++ {
		<-shop.BarbersDoneChan
	}

	close(shop.BarbersDoneChan)
	color.Green("--------------------------")
	color.Green("Shop is closed for the day")
}

func (shop *BarberShop) addClient(client string) {
	color.Green("%s arrives", client)

	if shop.Open {
		select {
		case shop.ClientsChan <- client:
			color.Yellow("%s is waiting for a barber", client)
		default:
			color.Red("%s is leaving because the shop is full", client)
		}
	} else {
		color.Red("Shop is closed so %s is sent home", client)
	}
}
