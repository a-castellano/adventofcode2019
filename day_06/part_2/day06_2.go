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

func findOrbiterPathToCOM(orbiterName string, orbiters map[string]*Orbit) []string {

	var orbitalPath []string

	var massName string

	massName = orbiters[orbiterName].mass.orbiter

	for massName != "COM" {
		orbitalPath = append(orbitalPath, massName)
		massName = orbiters[massName].mass.orbiter
	}
	orbitalPath = append(orbitalPath, "COM")

	return orbitalPath
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	orbiters := processFile(filename)

	pathFromYOUtoCOM := findOrbiterPathToCOM("YOU", orbiters)
	pathFromSANtoCOM := findOrbiterPathToCOM("SAN", orbiters)

	var found bool = false
	for youIndex, mass := range pathFromYOUtoCOM {
		for sanIndex, candidate := range pathFromSANtoCOM {
			if mass == candidate {
				fmt.Printf("Minimum number of orbital transfers required: %d\n", youIndex+sanIndex)
				found = true
				break
			}
		}
		if found {
			break
		}
	}
}
