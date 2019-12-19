// Álvaro Castellano Vela 2019/12/19
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
	cost    int
	up      *Point
	down    *Point
	left    *Point
	right   *Point
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

func explorePoint(input chan int, output chan int, X int, Y int, area map[string]*Point, goingValue int, returningValue int) (int, *Point) {

	var pointStr string
	var typeFound int
	newPoint := Point{X, Y, ' ', 0, nil, nil, nil, nil}
	pointStr = fmt.Sprintf("%d,%d", X, Y)

	if _, ok := area[pointStr]; !ok {
		input <- goingValue
		robotResponse := <-output
		typeFound = robotResponse
		switch robotResponse {
		case 0:
			newPoint.content = '#'
			newPoint.cost = 10000000
		case 1:
			newPoint.content = '.'
			newPoint.cost = 1
			input <- returningValue //  Return
			<-output                // We already know this output
		case 2:
			newPoint.content = 'O'
			newPoint.cost = 1
			input <- returningValue // South Return
			<-output                // We already know this output
		}
	} else {
		return -1, nil
	}
	return typeFound, &newPoint
}

func exploreArea(input chan int, output chan int, X int, Y int, area map[string]*Point) bool {

	var continueExploration bool = true

	var nextX int = X
	var nextY int = Y
	var pointStr string
	var nextPointStr string
	var pointsToExplore []*Point
	var nextPoint *Point
	var typeFound int
	pointStr = fmt.Sprintf("%d,%d", nextX, nextY)
	currentPoint := area[pointStr]

	//Explore up
	nextX--
	nextPointStr = fmt.Sprintf("%d,%d", nextX, nextY)
	typeFound, nextPoint = explorePoint(input, output, nextX, nextY, area, 1, 2)
	if typeFound != -1 {
		nextPoint.down = currentPoint
		currentPoint.up = nextPoint
		area[nextPointStr] = nextPoint
		if typeFound > 0 {
			pointsToExplore = append(pointsToExplore, nextPoint)
		}
	}
	nextX++

	//Explore Down
	nextX++
	nextPointStr = fmt.Sprintf("%d,%d", nextX, nextY)
	typeFound, nextPoint = explorePoint(input, output, nextX, nextY, area, 2, 1)
	if typeFound != -1 {
		nextPoint.up = currentPoint
		currentPoint.down = nextPoint
		area[nextPointStr] = nextPoint
		if typeFound > 0 {
			pointsToExplore = append(pointsToExplore, nextPoint)
		}
	}
	nextX--

	//Explore Left
	nextY--
	nextPointStr = fmt.Sprintf("%d,%d", nextX, nextY)
	typeFound, nextPoint = explorePoint(input, output, nextX, nextY, area, 3, 4)
	if typeFound != -1 {
		nextPoint.right = currentPoint
		currentPoint.left = nextPoint
		area[nextPointStr] = nextPoint
		if typeFound > 0 {
			pointsToExplore = append(pointsToExplore, nextPoint)
		}
	}
	nextY++

	//Explore Right
	nextY++
	nextPointStr = fmt.Sprintf("%d,%d", nextX, nextY)
	typeFound, nextPoint = explorePoint(input, output, nextX, nextY, area, 4, 3)
	if typeFound != -1 {
		nextPoint.left = currentPoint
		currentPoint.right = nextPoint
		area[nextPointStr] = nextPoint
		if typeFound > 0 {
			pointsToExplore = append(pointsToExplore, nextPoint)
		}
	}
	nextY--

	if len(pointsToExplore) == 0 {
		continueExploration = false
	} else {
		for _, pointToxplore := range pointsToExplore {
			var goingValue int
			var returnigValue int
			if pointToxplore.X == X-1 {
				goingValue = 1    //Up
				returnigValue = 2 //Down
			}
			if pointToxplore.X == X+1 {
				goingValue = 2    //Down
				returnigValue = 1 //Up
			}
			if pointToxplore.Y == Y-1 {
				goingValue = 3    //Left
				returnigValue = 4 //Right
			}
			if pointToxplore.Y == Y+1 {
				goingValue = 4    //Right
				returnigValue = 3 //Left
			}
			input <- goingValue
			<-output // We already know this output
			exploreArea(input, output, pointToxplore.X, pointToxplore.Y, area)
			input <- returnigValue
			<-output // We already know this output
		}
	}
	return continueExploration
}

func getArea(intCode []int) (map[string]*Point, *Point) {
	input := make(chan int)
	output := make(chan int)

	area := make(map[string]*Point)
	var X int = 0
	var Y int = 0
	var oxigenPoint *Point

	firstPoint := Point{X, Y, '.', 1, nil, nil, nil, nil}
	pointStr := fmt.Sprintf("%d,%d", X, Y)
	area[pointStr] = &firstPoint

	go runIntCode(intCode, input, output)

	exploreArea(input, output, X, Y, area)

	for _, point := range area {
		if point.content == 'O' {
			oxigenPoint = point
		}
	}

	return area, oxigenPoint
}

func fillOxigen(area map[string]*Point, oxigenPoint *Point) int {
	var continueFilling bool = true
	var nextPointsToFill []*Point
	var minutes = -1
	nextPointsToFill = append(nextPointsToFill, oxigenPoint)

	var newNextPointsToFill []*Point
	for continueFilling {
		for _, nextPoint := range nextPointsToFill {
			nextPoint.content = 'O'
			if nextPoint.up != nil {
				if nextPoint.up.content == '.' {
					newNextPointsToFill = append(newNextPointsToFill, nextPoint.up)
				}
			}
			if nextPoint.down != nil {
				if nextPoint.down.content == '.' {
					newNextPointsToFill = append(newNextPointsToFill, nextPoint.down)
				}
			}
			if nextPoint.right != nil {
				if nextPoint.right.content == '.' {
					newNextPointsToFill = append(newNextPointsToFill, nextPoint.right)
				}
			}
			if nextPoint.left != nil {
				if nextPoint.left.content == '.' {
					newNextPointsToFill = append(newNextPointsToFill, nextPoint.left)
				}
			}
		}
		nextPointsToFill = make([]*Point, 0)
		nextPointsToFill = newNextPointsToFill
		newNextPointsToFill = make([]*Point, 0)
		continueFilling = false
		for _, point := range area {
			if point.content == '.' {
				continueFilling = true
				break
			}
		}
		minutes++

	}

	return minutes
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	intCode := processFile(filename)
	area, oxigenPoint := getArea(intCode)
	minutes := fillOxigen(area, oxigenPoint)
	fmt.Println("Minutes until oxigen fills area: ", minutes)
}
