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

func processFile(filename string, desiredOutput int) int {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var originalIntCode []int
	var noun, verb int

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
		originalIntCode = append(originalIntCode, value)
	}

	// Before running the program, replace position 1 with the value 12 and replace position 2 with the value 2
	intCode := make([]int, len(originalIntCode))
	for i := 0; i <= 99; i++ {
		for j := 0; j <= 99; j++ {
			copy(intCode, originalIntCode)
			intCode[1] = i
			intCode[2] = j
			runIntCode(intCode)
			if intCode[0] == desiredOutput {
				noun = i
				verb = j
				return noun*100 + verb
			}
		}
	}
	return -1
}

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		log.Fatal("You must supply a file to process and desired output value.")
	}
	filename := args[0]
	output, _ := strconv.Atoi(args[1])
	result := processFile(filename, output)
	fmt.Printf("Result: %d\n", result)
}
