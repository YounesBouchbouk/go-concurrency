package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailer, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	syccess     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch

	return <-ch
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++

	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1

		fmt.Printf("received order #%d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1

		ms := ""

		success := false

		if rnd < 5 {
			pizzasFailer++
		} else {
			pizzasMade++
		}

		total++

		fmt.Printf("Making pizza #%d. It will take %d secondes .... \n", pizzaNumber, delay)

		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {

			ms = fmt.Sprintf("*** We ran out of ingredients for  pizza #%d", pizzaNumber)

		} else if rnd <= 4 {
			ms = fmt.Sprintf("*** The cook quir while making pizza #%d", pizzaNumber)

		} else {
			success = true

			ms = fmt.Sprintf("Pizza order #%d is ready!", pizzaNumber)
		}

		return &PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     ms,
			syccess:     success,
		}

	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}

}

func pizziria(pizzaMaker *Producer) {
	//keep track of with pizza we are making
	var i = 0
	//run forever o until receive a quit notification

	//try to mke pizzas

	for {
		// try to make pizza
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber

			select {
			case pizzaMaker.data <- *currentPizza:

			case quitChan := <-pizzaMaker.quit:
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}

		//decision
	}
}

func main() {
	//seed the random number generatot
	rand.Seed(time.Now().UnixNano())

	//print message

	color.Cyan("The Pizzeria is open for business !")
	color.Cyan("----------------------------------")

	// create a producer

	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	go pizziria(pizzaJob)

	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.syccess {
				color.Green(i.message)
				color.Green("Order #%d is out for delivery ", i.pizzaNumber)

			} else {
				color.Red(i.message)
				color.Red("the costumer is really mad! ", i.pizzaNumber)
			}

		} else {

			color.Cyan("the end of the day ....")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("some thing went wrong : ", err)
			}
		}
	}

	color.Cyan("--------------------------------------------")
	color.Yellow("we made %d pizzas , buit %d has been failed , and we have made %d pizzas successfully", pizzasMade, pizzasFailer, total)

}
