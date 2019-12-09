// Álvaro Castellano Vela 2019/12/09
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func decodeInstruction(instruction int) (int, int, int, int) {
	var opcode int = instruction % 100
	var modes int = instruction / 100
	var mode1 = modes % 10
	modes = modes / 10
	var mode2 = modes % 10
	modes = modes / 10
	var mode3 = modes % 10

	return opcode, mode1, mode2, mode3
}

func runIntCode(intCode []int, input int) {
	var stop bool = false

	var position int = 0
	var output int = 0

	var relativeBase int = 0

	var parameter1, parameter2 int = 0, 0
	for stop != true {
		parameter1, parameter2 = 0, 0
		opcode, mode1, mode2, mode3 := decodeInstruction(intCode[position])
		switch opcode {
		case 1, 2:
			switch mode1 {
			case 0:
				parameter1 = intCode[intCode[position+1]]
			case 1:
				parameter1 = intCode[position+1]
			case 2:
				parameter1 = intCode[intCode[position+1]+relativeBase]
			}
			switch mode2 {
			case 0:
				parameter2 = intCode[intCode[position+2]]
			case 1:
				parameter2 = intCode[position+2]
			case 2:
				parameter2 = intCode[intCode[position+2]+relativeBase]
			}
			switch mode3 {
			case 0:
				if opcode == 1 {
					intCode[intCode[position+3]] = parameter1 + parameter2
				} else {
					intCode[intCode[position+3]] = parameter1 * parameter2
				}
			case 2:
				if opcode == 1 {
					intCode[intCode[position+3]+relativeBase] = parameter1 + parameter2
				} else {
					intCode[intCode[position+3]+relativeBase] = parameter1 * parameter2
				}
			}
			position += 4
		case 3, 4:
			if opcode == 3 {
				switch mode1 {
				case 0:
					intCode[intCode[position+1]] = input
				case 1:
					intCode[position+1] = input
				case 2:
					intCode[intCode[position+1]+relativeBase] = input
				}
			} else {
				switch mode1 {
				case 0:
					output = intCode[intCode[position+1]]
				case 1:
					output = intCode[position+1]
				case 2:
					output = intCode[intCode[position+1]+relativeBase]
				}
				if output != 0 {
					fmt.Printf("Output %d\n", output)
				}
			}
			position += 2
		case 5:
			var parameter int
			switch mode1 {
			case 0:
				parameter = intCode[intCode[position+1]]
			case 1:
				parameter = intCode[position+1]
			case 2:
				parameter = intCode[intCode[position+1]+relativeBase]
			}
			if parameter != 0 {
				switch mode2 {
				case 0:
					position = intCode[intCode[position+2]]
				case 1:
					position = intCode[position+2]
				case 2:
					position = intCode[intCode[position+2]+relativeBase]
				}
			} else {
				position += 3
			}
		case 6:
			var parameter int
			switch mode1 {
			case 0:
				parameter = intCode[intCode[position+1]]
			case 1:
				parameter = intCode[position+1]
			case 2:
				parameter = intCode[intCode[position+1]+relativeBase]
			}
			if parameter == 0 {
				switch mode2 {
				case 0:
					position = intCode[intCode[position+2]]
				case 1:
					position = intCode[position+2]
				case 2:
					position = intCode[intCode[position+2]+relativeBase]
				}
			} else {
				position += 3
			}
		case 7:
			switch mode1 {
			case 0:
				parameter1 = intCode[intCode[position+1]]
			case 1:
				parameter1 = intCode[position+1]
			case 2:
				parameter1 = intCode[intCode[position+1]+relativeBase]
			}
			switch mode2 {
			case 0:
				parameter2 = intCode[intCode[position+2]]
			case 1:
				parameter2 = intCode[position+2]
			case 2:
				parameter2 = intCode[intCode[position+2]+relativeBase]
			}
			switch mode3 {
			case 0:
				if parameter1 < parameter2 {
					intCode[intCode[position+3]] = 1
				} else {
					intCode[intCode[position+3]] = 0
				}
			case 2:
				if parameter1 < parameter2 {
					intCode[intCode[position+3]+relativeBase] = 1
				} else {
					intCode[intCode[position+3]+relativeBase] = 0
				}
			}
			position += 4
		case 8:
			switch mode1 {
			case 0:
				parameter1 = intCode[intCode[position+1]]
			case 1:
				parameter1 = intCode[position+1]
			case 2:
				parameter1 = intCode[intCode[position+1]+relativeBase]
			}
			switch mode2 {
			case 0:
				parameter2 = intCode[intCode[position+2]]
			case 1:
				parameter2 = intCode[position+2]
			case 2:
				parameter2 = intCode[intCode[position+2]+relativeBase]
			}
			switch mode3 {
			case 0:
				if parameter1 == parameter2 {
					intCode[intCode[position+3]] = 1
				} else {
					intCode[intCode[position+3]] = 0
				}
			case 2:
				if parameter1 == parameter2 {
					intCode[intCode[position+3]+relativeBase] = 1
				} else {
					intCode[intCode[position+3]+relativeBase] = 0
				}
			}
			position += 4
		case 9:
			switch mode1 {
			case 0:
				relativeBase += intCode[intCode[position+1]]
			case 1:
				relativeBase += intCode[position+1]
			case 2:
				relativeBase += intCode[intCode[position+1]+relativeBase]
			}
			position += 2

		case 99:
			stop = true
		default:
			log.Fatal("Unknown opcode ", opcode)
		}
	}
}

func processFile(filename string, input int) {
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
	for i := 0; i < 6000; i++ {
		intCode = append(intCode, 0)
	}

	runIntCode(intCode, input)
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	processFile(filename, 2)
}
