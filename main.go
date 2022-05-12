package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	const FILE_NAME = "data.txt"
	const MAX_BUFFER = 128*1024
	const SLICE_VALUE = 50

	//GET LENGTHS
	nrOfLines, nrOfElements := getDataLengths(FILE_NAME, MAX_BUFFER)
	fmt.Printf("Nr of lines: {%d}; Nr of elements: {%d}\n", nrOfLines, nrOfElements)

	//GET DATA INTO AN ARRAY
	//Init slice
	dataMatrix := make([][]float64, nrOfLines, 128*1024)
	for i := range dataMatrix {
		dataMatrix[i] = make([]float64, nrOfElements)
	}

	getDataIntoSlice(dataMatrix, nrOfLines, nrOfElements, FILE_NAME, MAX_BUFFER)

	//No channels=========================================================
	start := time.Now()
	sumOfSquaresValue := sumOfSquares(dataMatrix, 0, nrOfLines, nrOfElements)
	fmt.Printf("Sum of squares: {%f}\n", sumOfSquaresValue)
	defer finishProfiler(start)
	//====================================================================

	//Multiple channels ==================================================
	//myChannels := make(chan float64)
	//start := time.Now()
	//distributeProcesses(dataMatrix, SLICE_VALUE, myChannels, nrOfElements)
	//defer finishProfiler(start)
	//
	////Retrieve values
	//var resultsPerLine = make([]float64, nrOfLines)
	//
	//for i:=0; i<SLICE_VALUE; i++ {
	//	resultsPerLine[i] = <- myChannels
	//
	//	//Debug Log -> sum for each line:
	//	fmt.Printf("Value of line {%d}: %f \n", i, resultsPerLine[i])
	//}
	////Debug Log -> Total Sum
	//var totalSum float64 = 0;
	//
	//for i := 0; i< len(resultsPerLine); i++ {
	//	totalSum += resultsPerLine[i];
	//}
	//fmt.Printf("Total Sum: %f \n" , totalSum)
	//====================================================================
}

func finishProfiler(start time.Time)  {
	elapsed := time.Since(start)
	log.Printf("Execution took %d", elapsed.Nanoseconds())
}

func getDataIntoSlice(dataMx [][]float64, nrOfLines int, nrOfElements int, fileName string, maxBufferRead int) {

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	buf := []byte{}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, maxBufferRead) 	//custom buffer - adjusted

	lineCounter := 0

	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, " ")

		var indexCounter int = 0;

		for _, value := range items {
			if parsedValue, err := strconv.ParseFloat(value, 64); err == nil {
				dataMx[lineCounter][indexCounter] = parsedValue
				indexCounter++;
			}
		}
		indexCounter = 0;
		lineCounter++
	}
}

func getDataLengths(fileName string, maxBufferRead int) (int, int) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	buf := []byte{}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, maxBufferRead) 	//custom buffer - adjusted

	indexCounter := 0;
	lineCounter := 0

	for scanner.Scan() {

		line := scanner.Text()
		items := strings.Split(line, " ")

		//Count the number of values of the first line
		if lineCounter == 0 {
			for _, value := range items {
				if _, err := strconv.ParseFloat(value, 64); err == nil {
					indexCounter++;
				}
			}
		}

		lineCounter++
	}

	return lineCounter, indexCounter
}

func sumOfSquares(dataMx [][]float64, startLine int, endLine int, nrOfElements int) (float64) {
	var sum float64 = 0
	for i := startLine; i < endLine; i++ {
		for j := 0; j < nrOfElements; j++ {
			sum += math.Sqrt(dataMx[i][j]);
		}
	}
	return sum
}

func distributeProcesses(dataMx [][]float64, sliceValue int, channel chan float64, nrOfElements int) {

	//sliceValue == number of processes
	modulo := len(dataMx) % sliceValue
	fmt.Printf("Numbers Array Length: {%d}; Slice Value: {%d}; Modulo: {%d} \n", len(dataMx), sliceValue, modulo)

	if modulo != 0 {
		fmt.Printf("Error - Array Length can't be sliced equally. Aborting...\n")
		return
	}

	for i:=0; i<sliceValue; i++ {
		go channelSumOfSquares(dataMx, channel, sliceValue, nrOfElements, i);
	}
}

func channelSumOfSquares(dataMx [][]float64, channel chan float64, sliceValue int, nrOfElements int, indexOfPortion int) {

	chunkSize := len(dataMx) / sliceValue;
	sum := sumOfSquares(dataMx, indexOfPortion * chunkSize, indexOfPortion * chunkSize + chunkSize, nrOfElements)
	channel <- sum
}
