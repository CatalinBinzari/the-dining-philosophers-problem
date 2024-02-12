package main

import (
	"fmt"
	"sync"
	"time"
)

type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

// slice of Philosopher type
var philosophers = []Philosopher{
	{name: "Plato", rightFork: 4, leftFork: 0},
	{name: "Socrates", rightFork: 0, leftFork: 1},
	{name: "Aristotel", rightFork: 1, leftFork: 2},
	{name: "Pascal", rightFork: 2, leftFork: 3},
	{name: "Locke", rightFork: 3, leftFork: 4},
}

// how many time can a philosopher eat before die
var hunger = 3

// eating time
var eatTime = 1 * time.Second

// after they finish eating they think some time
var thinkTime = 3 * time.Second

// pause time
var sleepTime = 1 * time.Second

// Philosopher is struct type that contains name, rightFork and leftFork. The name is a string type, rightFork and leftFork are int type
func main() {

	// print out a wlc message
	fmt.Println("dining philosophers problem")
	fmt.Println("--------------------------")
	fmt.Println("the table is empty")

	// start the meal
	dine()

	// print out a bye message
	fmt.Println("the table is empty")

}

// create 5 goroutine, one for each philosopher
func dine() {
	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	// forks is a map of all 5 forks
	// once we create a mutex, we should not copy it, but we can use a pointer to it
	var forks = make(map[int]*sync.Mutex)
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	var finishOrderMutex = &sync.Mutex{}
	var finishOrderList = make([]string, 0, len(philosophers))

	fmt.Println(len(finishOrderList), cap(finishOrderList))

	// start the meal.
	for i := 0; i < len(philosophers); i++ {
		// fire off a goroutine for the current philosopher
		go diningProblem(philosophers[i], wg, forks, seated, &finishOrderList, finishOrderMutex)
	}

	wg.Wait()

	// print the order of the philosophers finished eating
	fmt.Println("Finish order list:", finishOrderList)
	fmt.Println(len(finishOrderList), cap(finishOrderList))

}

func diningProblem(philosopher Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seated *sync.WaitGroup, finishOrderList *[]string, finishOrderMutex *sync.Mutex) {
	defer wg.Done()

	// seat the philosopher at thte table
	fmt.Printf("%s seating\n", philosopher.name)
	seated.Done()

	seated.Wait()
	fmt.Printf("%s seated\n", philosopher.name)

	// eat three times
	for i := hunger; i > 0; i-- {
		// needs to check fork from the left and fork from the right
		// get a lock on both forks
		forks[philosopher.leftFork].Lock() // if someone used has the fork, this goroutine pauses until the fork is released
		fmt.Printf("%s checks the left fork.\n", philosopher.name)
		forks[philosopher.rightFork].Lock()
		fmt.Printf("%s checks the right fork.\n", philosopher.name)

		fmt.Printf("--------%s has both forks, and start eating.(%d/%d)\n", philosopher.name, i, hunger)
		time.Sleep(eatTime)

		fmt.Printf("%s is thinking.\n", philosopher.name)
		time.Sleep(thinkTime)

		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()

		fmt.Printf("%s put down the fork.\n", philosopher.name)
	}

	fmt.Println(philosopher.name, "is satisfied and leaves the table.")
	finishOrderMutex.Lock()
	*finishOrderList = append(*finishOrderList, philosopher.name)
	fmt.Println(philosopher.name, "is added to the list.")
	finishOrderMutex.Unlock()

	fmt.Println(len(*finishOrderList), cap(*finishOrderList))
}
