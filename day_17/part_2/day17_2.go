// Álvaro Castellano Vela 2019/12/21
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	X       int
	Y       int
	content rune
	ID      int
	up      int
	down    int
	left    int
	right   int
}

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

func processFile(filename string) []int {
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

	return intCode
}

func runIntCode(intCode []int, input chan int, output chan int) {
	var stop bool = false

	var position int = 0

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
					intCode[intCode[position+1]] = <-input
				case 1:
					intCode[position+1] = <-input
				case 2:
					intCode[intCode[position+1]+relativeBase] = <-input
				}
			} else {
				switch mode1 {
				case 0:
					output <- intCode[intCode[position+1]]
				case 1:
					output <- intCode[position+1]
				case 2:
					output <- intCode[intCode[position+1]+relativeBase]
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
			fmt.Println("END")
			output <- 99
		default:
			log.Fatal("Unknown opcode ", opcode)
		}
	}
}

func isIntersection(shieldMap [][]rune, rows int, columns int, i int, j int) bool {
	var intersection bool = false

	if i > 0 && i < (rows-2) && j > 0 && j < (columns-2) {
		if shieldMap[i][j] == shieldMap[i-1][j] && shieldMap[i-1][j] == shieldMap[i+1][j] && shieldMap[i+1][j] == shieldMap[i][j-1] && shieldMap[i][j-1] == shieldMap[i][j+1] {
			intersection = true
		}
	}
	return intersection
}

func getShieldMap(intCode []int) ([][]rune, [][]*Point) {
	input := make(chan int)
	output := make(chan int)

	var shieldMap [][]rune
	var shieldPoints [][]*Point
	var shieldRow []rune
	var shieldPointsRow []*Point
	var cameraOutput int

	go runIntCode(intCode, input, output)

	cameraOutput = <-output

	for cameraOutput != 99 {
		if cameraOutput != 10 {
			shieldRow = append(shieldRow, rune(cameraOutput))

			newPoint := Point{-1, -1, '.', -1, -1, -1, -1, -1}
			shieldPointsRow = append(shieldPointsRow, &newPoint)
		} else {
			if len(shieldRow) != 0 {
				shieldMap = append(shieldMap, shieldRow)
				shieldPoints = append(shieldPoints, shieldPointsRow)
				shieldRow = nil
				shieldPointsRow = nil
			}
		}
		cameraOutput = <-output
	}
	return shieldMap, shieldPoints
}

func findPath(shieldMap [][]rune, shieldPoints [][]*Point) {
	var rows int = len(shieldMap)
	var columns int = len(shieldMap[0])

	var startPoint Point
	startPoint.up = -1
	startPoint.down = -1
	startPoint.left = -1
	startPoint.right = -1
	points := make(map[int]*Point)

	var ID int = 1

	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			fmt.Printf("%c", shieldMap[i][j])
			if shieldMap[i][j] == '#' {
				shieldPoints[i][j].ID = ID
				shieldPoints[i][j].X = i
				shieldPoints[i][j].Y = j
				shieldPoints[i][j].content = '#'
				points[ID] = shieldPoints[i][j]
				ID++
			} else {
				if shieldMap[i][j] != '.' {
					startPoint.X = i
					startPoint.Y = j
					startPoint.ID = 0
					startPoint.content = shieldMap[i][j]
					shieldPoints[i][j] = &startPoint
					points[0] = &startPoint
				}
			}
		}
		fmt.Printf("\n")
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {

			if shieldPoints[i][j].content != '.' {
				if i != 0 {
					if shieldPoints[i-1][j].content != '.' {
						shieldPoints[i][j].up = shieldPoints[i-1][j].ID
					}
				}
				if i != rows-1 {
					if i != rows-1 {
						shieldPoints[i][j].down = shieldPoints[i+1][j].ID
					}
				}
				if j != 0 {
					if shieldPoints[i][j-1].content != '.' {
						shieldPoints[i][j].left = shieldPoints[i][j-1].ID
					}
				}
				if j != columns-1 {
					if shieldPoints[i][j+1].content != '.' {
						shieldPoints[i][j].right = shieldPoints[i][j+1].ID
					}
				}

			}

		}
	}

	direction := 'L'
	switch startPoint.content {
	case '^':
		direction = 'U'
	case '>':
		direction = 'R'
	case '<':
		direction = 'L'
	case 'v':
		direction = 'D'
	}
	visitedPoints := make(map[int]bool)
	currentPosition := &startPoint
	visitedPoints[currentPosition.ID] = true
	var directions []rune
	for len(visitedPoints)+1 != len(points) {
		var nextStep int = -1
		//Try to follow the tile
		switch direction {
		case 'U':
			if _, ok := visitedPoints[currentPosition.up]; !ok && currentPosition.up != -1 {
				nextStep = currentPosition.up
			} else {
				if currentPosition.up != -1 {
					nextStep = currentPosition.up
				}
			}
		case 'R':
			if _, ok := visitedPoints[currentPosition.right]; !ok && currentPosition.right != -1 {
				nextStep = currentPosition.right
			} else {
				if currentPosition.right != -1 {
					nextStep = currentPosition.right
				}
			}
		case 'L':
			if _, ok := visitedPoints[currentPosition.left]; !ok && currentPosition.left != -1 {
				nextStep = currentPosition.left
			} else {
				if currentPosition.left != -1 {
					nextStep = currentPosition.left
				}
			}
		case 'D':
			if _, ok := visitedPoints[currentPosition.down]; !ok && currentPosition.down != -1 {
				nextStep = currentPosition.down
			} else {
				if currentPosition.down != -1 {
					nextStep = currentPosition.down
				}
			}
		}
		if nextStep == -1 {
			if _, ok := visitedPoints[currentPosition.up]; !ok && currentPosition.up != -1 {
				nextStep = currentPosition.up
				direction = 'U'
			}
			if _, ok := visitedPoints[currentPosition.right]; !ok && currentPosition.right != -1 {
				nextStep = currentPosition.right
				direction = 'R'
			}
			if _, ok := visitedPoints[currentPosition.left]; !ok && currentPosition.left != -1 {
				nextStep = currentPosition.left
				direction = 'L'
			}
			if _, ok := visitedPoints[currentPosition.down]; !ok && currentPosition.down != -1 {
				nextStep = currentPosition.down
				direction = 'D'
			}
		}
		visitedPoints[currentPosition.ID] = true
		shieldMap[currentPosition.X][currentPosition.Y] = '_'
		currentPosition = points[nextStep]
		directions = append(directions, direction)

	}
	fmt.Println("\nInstructions required for visiting every part of the scaffold at least once:")
	var lastDirection rune = '-'
	var turn rune = 'L'
	var stepCounter int = -1
	for _, direction := range directions {
		if direction != lastDirection {
			if stepCounter != -1 {
				fmt.Printf("%d ", stepCounter+1)
			}
			switch lastDirection {
			case '-':
				lastDirection = 'R'
			case 'U':
				switch direction {
				case 'U':
					turn = '?'
				case 'D':
					turn = '?'
				case 'L':
					turn = 'L'
				case 'R':
					turn = 'R'
				}
			case 'L':
				switch direction {
				case 'U':
					turn = 'R'
				case 'D':
					turn = 'L'
				case 'L':
					turn = '?'
				case 'R':
					turn = '?'
				}
			case 'R':
				switch direction {
				case 'U':
					turn = 'L'
				case 'D':
					turn = 'R'
				case 'L':
					turn = '?'
				case 'R':
					turn = '?'
				}
			case 'D':
				switch direction {
				case 'U':
					turn = '?'
				case 'D':
					turn = '?'
				case 'L':
					turn = 'R'
				case 'R':
					turn = 'L'
				}

			}
			fmt.Printf("%c ", turn)
			lastDirection = direction
			stepCounter = 0
		} else {
			stepCounter++
		}
	}
	fmt.Printf("%d ", stepCounter+1)
	fmt.Printf("\n")
}

func controlVacuum(intCode []int, A []rune, B []rune, C []rune, main_routine []rune) int {
	output := make(chan int, 1)

	var vacuumInput []int

	var result int
	for _, asci := range main_routine {
		vacuumInput = append(vacuumInput, int(asci))
	}
	vacuumInput = append(vacuumInput, 10)
	for _, asci := range A {
		vacuumInput = append(vacuumInput, int(asci))
	}
	vacuumInput = append(vacuumInput, 10)

	for _, asci := range B {
		vacuumInput = append(vacuumInput, int(asci))
	}
	vacuumInput = append(vacuumInput, 10)
	for _, asci := range C {
		vacuumInput = append(vacuumInput, int(asci))
	}
	vacuumInput = append(vacuumInput, 10)
	vacuumInput = append(vacuumInput, int('n'))
	vacuumInput = append(vacuumInput, 10)

	input := make(chan int, len(vacuumInput))
	for i := 0; i < len(vacuumInput); i++ {
		input <- vacuumInput[i]
	}
	go runIntCode(intCode, input, output)

	result = -1
	//	for true {
	//		result = <-output
	//		fmt.Println(result)
	//	}
	for result < 255 {
		result = <-output
	}

	return result
}

func main() {
	var A, B, C, main_routine []rune
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	intCode := processFile(filename)
	intCodeCopy := make([]int, len(intCode))
	copy(intCodeCopy, intCode)
	shieldMap, shieldPoints := getShieldMap(intCode)
	findPath(shieldMap, shieldPoints)

	A = []rune("L,12,L,12,R,12")
	B = []rune("L,8,L,8,R,12,L,8,L,8")
	C = []rune("L,10,R,8,R,12")
	main_routine = []rune("A,A,B,C,C,A,B,C,A,B")

	fmt.Println("A -> Size: ", len(A), " -> ", A)
	fmt.Println("B -> Size: ", len(B), " -> ", B)
	fmt.Println("C -> Size: ", len(C), " -> ", C)
	fmt.Println("main_routine -> Size: ", len(main_routine), " -> ", main_routine)

	//Force the vacuum robot to wake up by changing the value in your ASCII program at address 0 from 1 to 2.
	intCodeCopy[0] = 2
	result := controlVacuum(intCodeCopy, A, B, C, main_routine)
	fmt.Println("Result: ", result)
}
