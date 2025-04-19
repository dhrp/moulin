package main

import (
	"fmt"
	"time"
)

func main() {

	defer println("Main Function Returned Here")
	waitChan := make(chan bool)
	// stopHeartbeatChan := make(chan bool)

	var isWorking = struct{ active bool }{active: true}

	go MyWork(5, waitChan)
	go MySleep(15, &isWorking)

	<-waitChan
	isWorking.active = false

	time.Sleep(5 * time.Second) //wait to see

}

// MySleep sleeps for 1 second in each iteration
func MySleep(sleepTime int, isWorking *struct{ active bool }) {
	for i := 0; i < sleepTime; i++ {
		fmt.Println("goroutine running: ", i)
		time.Sleep(1 * time.Second)
		if isWorking.active == false {
			fmt.Println("done waiting")
			break
		}
	}
}

// MyWork does actual work
func MyWork(count time.Duration, c chan bool) {
	defer println("Work done")

	fmt.Println("Working! ")
	time.Sleep(count * time.Second)

	c <- true
}
