// Álvaro Castellano Vela 2019/12/02
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func runIntCode(intCode []int) {
	var stop bool = false

	var position int = 0

	for stop != true {
		if intCode[position] == 1 {
			intCode[intCode[position+3]] = intCode[intCode[position+1]] + intCode[intCode[position+2]]
		} else {
			if intCode[position] == 2 {
				intCode[intCode[position+3]] = intCode[intCode[position+1]] * intCode[intCode[position+2]]
			} else {
				if intCode[position] == 99 {
					stop = true
				} else {
					log.Fatal("Unknown opcode")
				}
			}
		}
		position += 4
	}
}

func processFile(filename string) int {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var intCode []int

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()

	line := scanner.Text()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	intCodeSlice := strings.Split(line, ",")

	for _, stringValue := range intCodeSlice {
		value, _ := strconv.Atoi(stringValue)
		intCode = append(intCode, value)
	}

	// Before running the program, replace position 1 with the value 12 and replace position 2 with the value 2
	intCode[1] = 12
	intCode[2] = 2
	runIntCode(intCode)

	return intCode[0]
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	result := processFile(filename)
	fmt.Printf("Position 0 has value: %d\n", result)
}
