package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var seatingCapacity = 5
var arrivalRate = 100
var cutDuraction = 1000 * time.Millisecond
var timeOpen = 5 * time.Second

func main() {

	rand.Seed(time.Now().UnixNano())

	color.Yellow("the sleeping Barber Prolem ")
	color.Yellow("-----------------------------")

	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	barberShop := BarberShop{
		numOfbabrbers:  0,
		cutDuraction:   cutDuraction,
		IsOpen:         true,
		SeatCapacity:   seatingCapacity,
		ClientChan:     clientChan,
		BarberDoneChan: doneChan,
	}

	barberShop.addBarber("younes")
	barberShop.addBarber("Ommar")

	// start the barber and close if after timeOpen
	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(timeOpen)
		shopClosing <- true

		barberShop.closeshopForDay()

		closed <- true

	}()

	// add clients

	i := 0

	go func() {

		for {
			randomMillsecods := rand.Int() % (2 * arrivalRate)
			select {
			case <-time.After(time.Microsecond * time.Duration(randomMillsecods)):
				barberShop.addNewClient(fmt.Sprintf("***clientN#%d", i))
				i++
			case <-shopClosing:
				return
			}
		}

	}()

	<-closed

	// time.Sleep(3 * time.Second)
}
