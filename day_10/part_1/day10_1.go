// Álvaro Castellano Vela 2019/12/10
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Point struct {
	X int
	Y int
}

type Asteroid struct {
	X              int
	Y              int
	FoundAsteroids int
}

func processFile(filename string) ([][]rune, [][]*Asteroid, []*Asteroid) {

	var space [][]rune
	var asteroidsMatrix [][]*Asteroid
	var asteroids []*Asteroid

	var rows int = 0
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "")
		space_line := make([]rune, len(line))
		asteroids_line := make([]*Asteroid, len(line))
		for position, stringvalue := range line {
			value := rune(stringvalue[0])
			space_line[position] = value
			if value == '#' {
				asteroid := Asteroid{rows, position, 0}
				asteroids_line[position] = &asteroid
				asteroids = append(asteroids, &asteroid)
			}
		}
		space = append(space, space_line)
		asteroidsMatrix = append(asteroidsMatrix, asteroids_line)
		rows++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return space, asteroidsMatrix, asteroids

}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func realpoint(point int) int {
	return point*1024 + 512
}

func bresenham(x0, y0, x1, y1 int) []Point {

	var points []Point

	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y0
	if dy < 0 {
		dy = -dy
	}
	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy

	for {
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
		points = append(points, Point{x0, y0})
	}
	return points[:len(points)-1]
}

func calculateIntersection(x0, y0, x1, y1 int) float64 {

	rx0 := realpoint(x0)
	ry0 := realpoint(y0)
	rx1 := realpoint(x1)
	ry1 := realpoint(y1)

	tan := (float64(Abs(ry1-ry0)) / float64(Abs(rx1-rx0)))

	return tan

}

func findAsteroids(space [][]rune, asteroidsMatrix [][]*Asteroid, asteroids []*Asteroid) int {

	for asteroidID, asteroid := range asteroids {
		for targetID, targetAsteroid := range asteroids {
			if asteroidID != targetID {
				points := bresenham(asteroid.X, asteroid.Y, targetAsteroid.X, targetAsteroid.Y)

				tans_used := make(map[float64]bool)
				tan_to_target := calculateIntersection(asteroid.X, asteroid.Y, targetAsteroid.X, targetAsteroid.Y)
				for _, point := range points {
					if asteroidsMatrix[point.X][point.Y] != nil {
						tan := calculateIntersection(asteroid.X, asteroid.Y, point.X, point.Y)
						tans_used[tan] = true
					}
				}
				if _, ok := tans_used[tan_to_target]; !ok {
					asteroid.FoundAsteroids++
				}

			}
		}
	}

	var maxFound int = 0

	for _, asteroid := range asteroids {
		if asteroid.FoundAsteroids > maxFound {
			maxFound = asteroid.FoundAsteroids
		}
	}
	return maxFound
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	space, asteroidsMatrix, asteroids := processFile(filename)

	result := findAsteroids(space, asteroidsMatrix, asteroids)
	fmt.Println("Result: ", result)
}
