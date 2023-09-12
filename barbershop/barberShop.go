package main

import (
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	numOfbabrbers  int
	SeatCapacity   int
	cutDuraction   time.Duration
	BarberDoneChan chan bool
	ClientChan     chan string
	IsOpen         bool
}

func (shop *BarberShop) addBarber(name string) {

	shop.numOfbabrbers++

	go func() {

		isSleeping := false
		color.Green("the barber %s is enter the shop ", name)

		for {
			if len(shop.ClientChan) == 0 {
				color.Yellow("the barber %s didn't find any one ", name)
				color.Yellow("the barber %s is going to sleep now  ", name)
				isSleeping = true
			}

			clientname, shopOpen := <-shop.ClientChan

			if shopOpen {
				if isSleeping {
					color.Green("the barber %s wakes up ", name)
					isSleeping = false
				}
				shop.startHairCut(name, clientname)
			} else {
				shop.sendBarberHome()
				return
			}
		}

	}()
}

func (shop *BarberShop) startHairCut(babrbername, client string) {
	color.Yellow("the barber %s is going to start a haircut for   ", babrbername, client)
	time.Sleep(shop.cutDuraction)
	color.Yellow("the barber %s is finishing the haircut to %s   ", babrbername, client)

}

func (shop *BarberShop) sendBarberHome() {
	color.Yellow("the barber is going home")
	shop.BarberDoneChan <- true

}

func (shop *BarberShop) closeshopForDay() {
	color.Yellow("the barber is giong to be closed")

	close(shop.ClientChan)

	for i := 0; i < shop.numOfbabrbers; i++ {
		<-shop.BarberDoneChan
	}
	close(shop.BarberDoneChan)

	color.Green("----------------------------------------------")
	color.Green("the barbershop is closed for the day ")

}

func (shop *BarberShop) addNewClient(clientName string) {
	if shop.IsOpen {
		select {
		case shop.ClientChan <- clientName:
			color.Green("%s has tekes a seat in the waiting room", clientName)
		default:
			color.Red("The waiting room is full , so %s leaves ", clientName)
		}
	} else {
		color.Red("the client %s is leaving the barber because it's closes", clientName)
	}
}
