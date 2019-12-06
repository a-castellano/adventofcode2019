// Álvaro Castellano Vela 2019/12/06
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

type Orbit struct {
	orbiter string
	mass    *Orbit
}

func processFile(filename string) map[string]*Orbit {

	orbiters := make(map[string]*Orbit)

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		re := regexp.MustCompile("([[:alnum:]]+)\\)([[:alnum:]]+)$")
		match := re.FindAllStringSubmatch(scanner.Text(), -1)
		mass := match[0][1]
		orbiter := match[0][2]
		if _, ok := orbiters[mass]; !ok {
			newOrbit := Orbit{mass, nil}
			orbiters[mass] = &newOrbit
		}
		if _, ok := orbiters[orbiter]; ok {
			orbiters[orbiter].mass = orbiters[mass]
		} else {
			newOrbit := Orbit{orbiter, orbiters[mass]}
			orbiters[orbiter] = &newOrbit
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return orbiters
}

func countOrbits(orbiterName string, orbiters map[string]*Orbit, orbitersOrbits map[string]int) int {
	if orbiters[orbiterName].mass == nil {
		return 0
	}
	mass := orbiters[orbiterName].mass

	if orbits, ok := orbitersOrbits[mass.orbiter]; ok {
		return orbits + 1
	} else {
		orbitersOrbits[mass.orbiter] = countOrbits(mass.orbiter, orbiters, orbitersOrbits)
		return 1 + orbitersOrbits[mass.orbiter]
	}
}

func countTotalOrbits(orbiters map[string]*Orbit) int {

	orbitersOrbits := make(map[string]int)
	var totalOrbits int = 0

	for orbiterName, _ := range orbiters {
		orbitersOrbits[orbiterName] = countOrbits(orbiterName, orbiters, orbitersOrbits)
	}
	for _, orbits := range orbitersOrbits {
		totalOrbits += orbits
	}

	return totalOrbits
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	orbiters := processFile(filename)
	orbits := countTotalOrbits(orbiters)

	fmt.Printf("Total orbits: %d\n", orbits)
}
