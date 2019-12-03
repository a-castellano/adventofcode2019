// Álvaro Castellano Vela 2019/12/03
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Point struct {
	X int
	Y int
}

type WireDirection struct {
	turn  rune
	steps int
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func processFile(filename string) [][]WireDirection {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var directions [][]WireDirection

	directionRe := regexp.MustCompile("(R|U|L|D|)([[:digit:]]+)")

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		wireDirectionsText := strings.Split(scanner.Text(), ",")
		var wireDirections []WireDirection
		for _, direction := range wireDirectionsText {
			match := directionRe.FindAllStringSubmatch(direction, -1)
			turn := []rune(match[0][1])[0]
			steps, _ := strconv.Atoi(match[0][2])
			wireDirection := WireDirection{turn, steps}
			wireDirections = append(wireDirections, wireDirection)
		}
		directions = append(directions, wireDirections)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return directions
}

func getLoweLeftPoint(wiresDirections [][]WireDirection) int {

	matrix := [40000][40000]uint8{}
	var intersections []Point

	for i := 0; i < 40000; i++ {
		for j := 0; j < 40000; j++ {
			matrix[i][j] = 0
		}
	}

	for _, wiredirection := range wiresDirections {
		position := Point{20000, 20000}
		for _, direction := range wiredirection {
			switch direction.turn {
			case 'U':
				for i := 0; i < direction.steps; i++ {
					position.X--
					if matrix[position.X][position.Y] == 0 {
						matrix[position.X][position.Y] = 1
					} else {
						matrix[position.X][position.Y] = 2
						intersections = append(intersections, position)
					}
				}
			case 'D':
				for i := 0; i < direction.steps; i++ {
					position.X++
					if matrix[position.X][position.Y] == 0 {
						matrix[position.X][position.Y] = 1
					} else {
						matrix[position.X][position.Y] = 2
						intersections = append(intersections, position)
					}
				}
			case 'L':
				for i := 0; i < direction.steps; i++ {
					position.Y--
					if matrix[position.X][position.Y] == 0 {
						matrix[position.X][position.Y] = 1
					} else {
						matrix[position.X][position.Y] = 2
						intersections = append(intersections, position)
					}
				}
			case 'R':
				for i := 0; i < direction.steps; i++ {
					position.Y++
					if matrix[position.X][position.Y] == 0 {
						matrix[position.X][position.Y] = 1
					} else {
						matrix[position.X][position.Y] = 2
						intersections = append(intersections, position)
					}
				}
			}
		}
	}

	// Calculate minimum distance
	closestPoint := Point{-1, -1}
	closestPointDistance := Point{90000, 90000}
	for _, intersecton := range intersections {
		if Abs(intersecton.X-20000) <= closestPointDistance.X {
			if Abs(intersecton.Y-20000) <= closestPointDistance.Y {
				closestPoint = intersecton
				closestPointDistance.X = Abs(intersecton.X - 20000)
				closestPointDistance.Y = Abs(intersecton.Y - 20000)
			}
		}
	}

	minDistance := Abs(closestPoint.X-20000) + Abs(closestPoint.Y-20000)
	return minDistance
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	directions := processFile(filename)
	lowerLeftPoint := getLoweLeftPoint(directions)
	fmt.Printf("Lower-left Point distance: %d\n", lowerLeftPoint)
}
