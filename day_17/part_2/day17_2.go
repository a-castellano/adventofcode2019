// Álvaro Castellano Vela 2019/12/21
package main

import (
	"bufio"
	"fmt"
	"github.com/yourbasic/graph"
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
			output <- 99
		default:
			log.Fatal("Unknown opcode ", opcode)
		}
	}
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

func findPath(shieldMap [][]rune, shieldPoints [][]*Point) int {
	var result int

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
					if shieldPoints[i+1][j].content != '.' {
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

	gm := graph.New(len(points))
	for _, value := range points {
		if value.up != -1 {
			gm.Add(value.ID, points[value.up].ID)
		}
		if value.down != -1 {
			gm.Add(value.ID, points[value.down].ID)
		}
		if value.right != -1 {
			gm.Add(value.ID, points[value.right].ID)
		}
		if value.left != -1 {
			gm.Add(value.ID, points[value.left].ID)
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
	dist := make([]int, gm.Order())
	var counter int = 0
	graph.BFS(gm, 0, func(v, w int, _ int64) {
		changeDirection := false
		switch direction {
		case 'U':
			if points[v].X != points[w].X-1 {
				changeDirection = true
				//		fmt.Println("Change direction")
			}
		case 'D':
			if points[v].X != points[w].X+1 {
				changeDirection = true
				//		fmt.Println("Change direction")
			}
		case 'R':
			if points[v].Y != points[w].Y+1 {
				changeDirection = true
				//		fmt.Println("Change direction")
			}
		case 'L':
			if points[v].Y != points[w].Y-1 {
				changeDirection = true
				//		fmt.Println("Change direction")
			}
		}

		if changeDirection {
			//Determine new direction
			if points[v].X == points[w].X-1 {
				direction = 'U'
			}
			if points[v].X == points[w].X+1 {
				direction = 'D'
			}
			if points[v].Y == points[w].Y-1 {
				direction = 'L'
			}
			if points[v].Y == points[w].Y+1 {
				direction = 'R'
			}
			if counter > 0 {
				fmt.Printf("%d ", counter)
			}
			fmt.Printf("%c ", direction)
			counter = 0
		} else {
			counter++
		}

		//fmt.Printf("Current Direction: %c\n", direction)
		//fmt.Println(points[v].X, points[v].Y, "to", points[w].X, points[w].Y)
		dist[w] = dist[v] + 1
	})
	fmt.Printf("\n")
	return result
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	intCode := processFile(filename)
	shieldMap, shieldPoints := getShieldMap(intCode)
	result := findPath(shieldMap, shieldPoints)
	fmt.Println("Result: ", result)
}
