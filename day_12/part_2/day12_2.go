//ng  Álvaro Castellano Vela 2019/12/12
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

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

func calculateInitialStept(moons []Moon) int {

	pervious_states0 := make(map[string]bool)
	pervious_states1 := make(map[string]bool)
	pervious_states2 := make(map[string]bool)

	var exit0 int = -1
	var exit1 int = -1
	var exit2 int = -1

	for step := 0; true; step++ {
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
		}
		state0 := fmt.Sprintf("%d %d %d %d %d %d %d %d", moons[0].position.X, moons[1].position.X, moons[2].position.X, moons[3].position.X, moons[0].velocity.X, moons[1].velocity.X, moons[2].velocity.X, moons[3].velocity.X)
		state1 := fmt.Sprintf("%d %d %d %d %d %d %d %d", moons[0].position.Y, moons[1].position.Y, moons[2].position.Y, moons[3].position.Y, moons[0].velocity.Y, moons[1].velocity.Y, moons[2].velocity.Y, moons[3].velocity.Y)
		state2 := fmt.Sprintf("%d %d %d %d %d %d %d %d", moons[0].position.Z, moons[1].position.Z, moons[2].position.Z, moons[3].position.Z, moons[0].velocity.Z, moons[1].velocity.Z, moons[2].velocity.Z, moons[3].velocity.Z)
		if exit0 < 0 {
			if _, ok := pervious_states0[state0]; !ok {
				pervious_states0[state0] = true
			} else {
				exit0 = step
			}
		}

		if exit1 < 0 {
			if _, ok := pervious_states1[state1]; !ok {
				pervious_states1[state1] = true
			} else {
				exit1 = step
			}
		}

		if exit2 < 0 {
			if _, ok := pervious_states2[state2]; !ok {
				pervious_states2[state2] = true
			} else {
				exit2 = step
			}
		}

		if exit0 >= 0 && exit1 >= 0 && exit2 >= 0 {
			break
		}

	}
	return LCM(exit0, exit1, exit2)
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]

	moons := processFile(filename)

	initialStept := calculateInitialStept(moons)
	fmt.Println("Initial stept: ", initialStept)
}
