// Álvaro Castellano Vela 2019/12/01
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func calculateFuel(mass int) int {
	var fuel int
	fuel = mass/3 - 2
	if fuel <= 0 {
		fuel = 0
	} else {
		fuel += (calculateFuel(fuel))
	}
	return fuel
}

func processFile(filename string) int {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var total_module int = 0

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		mass, _ := strconv.Atoi(scanner.Text())
		total_module += calculateFuel(mass)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return total_module
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	total_module := processFile(filename)
	fmt.Printf("Total fuel needed: %d\n", total_module)
}
