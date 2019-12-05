// Álvaro Castellano Vela 2019/12/05
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

	var parameter1, parameter2 int = 0, 0
	for stop != true {
		parameter1, parameter2 = 0, 0
		opcode, mode1, mode2, mode3 := decodeInstruction(intCode[position])
		switch opcode {
		case 1, 2:
			if mode1 == 0 {
				parameter1 = intCode[intCode[position+1]]
			} else {
				parameter1 = intCode[position+1]
			}
			if mode2 == 0 {
				parameter2 = intCode[intCode[position+2]]
			} else {
				parameter2 = intCode[position+2]
			}
			if mode3 == 1 {
				//error
			} else {
				if opcode == 1 {
					intCode[intCode[position+3]] = parameter1 + parameter2
				} else {
					intCode[intCode[position+3]] = parameter1 * parameter2
				}
			}
			position += 4
		case 3, 4:
			if opcode == 3 {
				if mode1 == 0 {
					intCode[intCode[position+1]] = input
				} else {
					//error
				}
			} else {
				if mode1 == 0 {
					output = intCode[intCode[position+1]]
					if output != 0 {
						fmt.Printf("%d", output)
					}
				} else {
					//error
				}
			}
			position += 2
		case 5:
			var parameter int
			if mode1 == 0 {
				parameter = intCode[intCode[position+1]]
			} else {
				parameter = intCode[position+1]
			}
			if parameter != 0 {
				if mode2 == 0 {
					position = intCode[intCode[position+2]]
				} else {
					position = intCode[position+2]
				}
			} else {
				position += 3
			}
		case 6:
			var parameter int
			if mode1 == 0 {
				parameter = intCode[intCode[position+1]]
			} else {
				parameter = intCode[position+1]
			}
			if parameter == 0 {
				if mode2 == 0 {
					position = intCode[intCode[position+2]]
				} else {
					position = intCode[position+2]
				}
			} else {
				position += 3
			}
		case 7:
			if mode1 == 0 {
				parameter1 = intCode[intCode[position+1]]
			} else {
				parameter1 = intCode[position+1]
			}
			if mode2 == 0 {
				parameter2 = intCode[intCode[position+2]]
			} else {
				parameter2 = intCode[position+2]
			}
			if mode3 == 1 {
				//error
			} else {
				if parameter1 < parameter2 {
					intCode[intCode[position+3]] = 1
				} else {
					intCode[intCode[position+3]] = 0
				}
			}
			position += 4
		case 8:
			if mode1 == 0 {
				parameter1 = intCode[intCode[position+1]]
			} else {
				parameter1 = intCode[position+1]
			}
			if mode2 == 0 {
				parameter2 = intCode[intCode[position+2]]
			} else {
				parameter2 = intCode[position+2]
			}
			if mode3 == 1 {
				//error
			} else {
				if parameter1 == parameter2 {
					intCode[intCode[position+3]] = 1
				} else {
					intCode[intCode[position+3]] = 0
				}
			}
			position += 4

		case 99:
			stop = true
			fmt.Printf("\n")
		default:
			fmt.Printf("\n")
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

	runIntCode(intCode, input)
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	processFile(filename, 5)
}
