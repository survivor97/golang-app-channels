package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var FILE_NAME = "data.txt"

	//GET LENGTHS
	nrOfLines, nrOfElements := getDataLengths(FILE_NAME, 128*1024)
	fmt.Printf("Nr of lines: {%d}; Nr of elements: {%d}\n", nrOfLines, nrOfElements)

	//READ DATA INTO AN ARRAY
	//Init slice
	dataMatrix := make([][]float64, nrOfLines, 128*1024)
	for i := range dataMatrix {
		dataMatrix[i] = make([]float64, nrOfElements)
	}

	getDataIntoSlice(dataMatrix, nrOfLines, nrOfElements, FILE_NAME, 128*1024)

	fmt.Println(dataMatrix[99][0]);
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
