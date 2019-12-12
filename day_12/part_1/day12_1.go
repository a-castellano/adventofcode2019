// Álvaro Castellano Vela 2019/12/12
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Point struct {
	X int
	Y int
	Z int
}

type Velocity struct {
	X int
	Y int
	Z int
}

type Moon struct {
	position        Point
	velocity        Velocity
	potentialEnergy int
	kineticEnergy   int
}

func (moon *Moon) calculatePotentialEnergy() {
	moon.potentialEnergy = Abs(moon.position.X) + Abs(moon.position.Y) + Abs(moon.position.Z)
}

func (moon *Moon) calculateKineticEnergy() {
	moon.kineticEnergy = Abs(moon.velocity.X) + Abs(moon.velocity.Y) + Abs(moon.velocity.Z)
}

func (moon Moon) totalEnergy() int {
	return moon.potentialEnergy * moon.kineticEnergy
}

func (moon *Moon) applyVelocity() {
	moon.position.X += moon.velocity.X
	moon.position.Y += moon.velocity.Y
	moon.position.Z += moon.velocity.Z
}

func processFile(filename string) []Moon {

	var moons []Moon

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		re := regexp.MustCompile("<x=([-]?[[:digit:]]+), y=([-]?[[:digit:]]+), z=([-]?[[:digit:]]+)>$")
		match := re.FindAllStringSubmatch(scanner.Text(), -1)
		x, _ := strconv.Atoi(match[0][1])
		y, _ := strconv.Atoi(match[0][2])
		z, _ := strconv.Atoi(match[0][3])
		point := Point{x, y, z}
		velocity := Velocity{0, 0, 0}

		moon := Moon{point, velocity, 0, 0}
		moon.calculatePotentialEnergy()
		moon.calculateKineticEnergy()
		moons = append(moons, moon)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return moons
}

func calculateTotalEnergy(moons []Moon, steps int) int {

	var totalEnergy int = 0

	for step := 0; step < steps; step++ {
		for moonA, _ := range moons {
			for moonB, _ := range moons {
				if moonA != moonB {
					// Calculate Velocity
					positionA := moons[moonA].position
					positionB := moons[moonB].position
					if positionA.X > positionB.X {
						moons[moonA].velocity.X -= 1
					} else {
						if positionA.X < positionB.X {
							moons[moonA].velocity.X += 1
						}
					}

					if positionA.Y > positionB.Y {
						moons[moonA].velocity.Y -= 1
					} else {
						if positionA.Y < positionB.Y {
							moons[moonA].velocity.Y += 1
						}
					}

					if positionA.Z > positionB.Z {
						moons[moonA].velocity.Z -= 1
					} else {
						if positionA.Z < positionB.Z {
							moons[moonA].velocity.Z += 1
						}
					}
				}
			}
		}
		for moon, _ := range moons {
			moons[moon].applyVelocity()
			moons[moon].calculatePotentialEnergy()
			moons[moon].calculateKineticEnergy()
		}
	}
	for _, moon := range moons {
		totalEnergy += moon.totalEnergy()
	}
	return totalEnergy
}

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		log.Fatal("You must supply a file to process and the number of steps to calculate.")
	}
	filename := args[0]
	steps, _ := strconv.Atoi(args[1])

	moons := processFile(filename)

	totalEnergy := calculateTotalEnergy(moons, steps)
	fmt.Println("Total energy: ", totalEnergy)
}
