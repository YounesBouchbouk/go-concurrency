package main

import (
	"fmt"
	"strings"
)

func work(ping chan string, pong chan string) {

	for {
		val := <-ping

		pong <- fmt.Sprintf("%s !!!!", strings.ToUpper(val))
	}
}

func main() {

	ping := make(chan string)
	pong := make(chan string)

	for {

		var inpt string
		go work(ping, pong)

		fmt.Println("-> ")

		fmt.Scanln(&inpt)

		if inpt == strings.ToLower("q") {
			break
		} else {
			ping <- inpt
		}

		fmt.Printf("%s", <-pong)
	}
	close(ping)
	close(pong)

}
