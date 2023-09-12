package main

import (
	"fmt"
	"sync"
	"time"
)

type Philosopher struct {
	name      string
	leftFork  int
	rightFork int
}

var philosophers = []Philosopher{
	{name: "Plato", leftFork: 4, rightFork: 0},
	{name: "Socrat", leftFork: 0, rightFork: 1},
	{name: "Aristotle", leftFork: 1, rightFork: 2},
	{name: "Pascal", leftFork: 2, rightFork: 3},
	{name: "Locke", leftFork: 3, rightFork: 4},
}

var hunger = 3 // times philosopher ezat

var eatTime = 1 * time.Second

var thinkTime = 2 * time.Second

var sleepTime = 1 * time.Second

func main() {
	fmt.Println("Dinning Problem ")

	dine()

	fmt.Println("The table is empty.")

}

func dine() {

	// wg is the WaitGroup that keeps track of how many philosophers are still at the table. When
	// it reaches zero, everyone is finished eating and has left. We add 5 (the number of philosophers) to this
	// wait group.

	var wg = &sync.WaitGroup{}
	wg.Add(len(philosophers))

	// We want everyone to be seated before they start eating, so create a WaitGroup for that, and set it to 5.
	var seatedWg = &sync.WaitGroup{} // to wait intil all philosophers is on the table
	seatedWg.Add(len(philosophers))

	// forks is a map of all 5 forks. Forks are assigned using the fields leftFork and rightFork in the Philosopher
	// type. Each fork, then, can be found using the index (an integer), and each fork has a unique mutex.
	var forks = make(map[int]*sync.Mutex)

	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	// Start the meal by iterating through our slice of Philosophers.
	for i := 0; i < len(philosophers); i++ {
		go diningProblem(philosophers[i], wg, forks, seatedWg)
	}

	fmt.Println("The table is emty")

	// Wait for the philosophers to finish. This blocks until the wait group is 0.
	wg.Wait()

}

// diningProblem is the function fired off as a goroutine for each of our philosophers. It takes one
// philosopher, our WaitGroup to determine when everyone is done, a map containing the mutexes for every
// fork on the table, and a WaitGroup used to pause execution of every instance of this goroutine
// until everyone is seated at the table.
func diningProblem(philosopher Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seated *sync.WaitGroup) {

	defer wg.Done()

	fmt.Printf("%s is seated at the table.\n", philosopher.name)

	// Decrement the seated WaitGroup by one.
	seated.Done()

	// Wait until everyone is seated.
	seated.Wait()

	// Have this philosopher eatTime and thinkTime "hunger" times (3).
	for i := hunger; i > 0; i-- {
		// Get a lock on the left and right forks. We have to choose the lower numbered fork first in order
		// to avoid a logical race condition, which is not detected by the -race flag in tests; if we don't do this,
		// we have the potential for a deadlock, since two philosophers will wait endlessly for the same fork.
		// Note that the goroutine will block (pause) until it gets a lock on both the right and left forks.

		if philosopher.leftFork > philosopher.rightFork {
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the right fork.\n", philosopher.name)
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the left fork.\n", philosopher.name)
		} else {
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the left fork.\n", philosopher.name)
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the right fork.\n", philosopher.name)
		}

		// Byx the time we get to this line, the philosopher has a lock (mutex) on both forks.
		fmt.Printf("\t%s has both forks and is eating.\n", philosopher.name)
		time.Sleep(eatTime)

		// The philosopher starts to think, but does not drop the forks yet.
		fmt.Printf("\t%s is thinking.\n", philosopher.name)
		time.Sleep(thinkTime)

		// Unlock the mutexes for both forks.
		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()

		fmt.Printf("\t%s put down the forks.\n", philosopher.name)
	}

	// The philosopher has finished eating, so print out a message.
	fmt.Println(philosopher.name, "is satisified.")
	fmt.Println(philosopher.name, "left the table.")

}
