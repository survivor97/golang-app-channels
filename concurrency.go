package main

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"
)

var waitGroup = sync.WaitGroup{}

func main() {

	//Init
	var numbers [1000000]float64

	for i:=0; i<len(numbers); i++ {
		numbers[i] = float64(i)
	}

	start := time.Now()

	//Normal

	//additionOfSquares(numbers[:])

	//----------------------
	//Channels
	myChannels := make(chan float64)
	distributeProcesses(numbers[:],5,myChannels)

	for message := range myChannels {
		fmt.Printf("Channel val: %f \n", message)
	}

	//----------------------

	elapsed := time.Since(start)
	log.Printf("Function took %s", elapsed)
}

func additionOfSquares(numbersArray []float64) {
	//fmt.Printf("length: %d", len(numbersArray))
	var sum float64 = 0
	for i := 0; i < len(numbersArray); i++ {
		sum += math.Sqrt(numbersArray[i])
	}
	fmt.Printf("Sum value: %f \n", sum)
}

func distributeProcesses(numbersArray []float64, sliceValue int, channel chan float64) {

	//sliceValue == number of processes
	modulo := len(numbersArray) % sliceValue
	fmt.Printf("Numbers Array Length: {%d}; Slice Value: {%d}; Modulo: {%d} \n", len(numbersArray), sliceValue, modulo)

	if modulo != 0 {
		fmt.Printf("Error - Array Length can't be sliced equally. Aborting...\n")
		return
	}

	//go channelAdditionOfSquares(numbersArrays[:len(numbersArrays)/2], channel);

}

func channelAdditionOfSquares(numbersArray []float64, channel chan float64) {
	var sum float64 = 0
	for i := 0; i < len(numbersArray); i++ {
		sum += math.Sqrt(numbersArray[i])
	}

	channel <- sum

	close(channel)
	fmt.Printf("Sum value: %f \n", sum)
}

func valFunc(c chan float64, value float64) {
	var val float64 = value
	c<-val
}

func saySomething() {
	fmt.Println("Something")
	waitGroup.Done()
}

func other() {
	fmt.Println("Other")
	waitGroup.Done()
}
