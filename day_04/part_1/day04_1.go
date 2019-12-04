// Álvaro Castellano Vela 2019/12/04
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func meetCriteria(candidate int) bool {
	candidate_str := strconv.Itoa(candidate)
	if len(candidate_str) != 6 {
		return false
	} else {
		var numberBefore rune = 1
		sameAdjacent := false
		for _, value := range candidate_str {
			if value < numberBefore {
				return false
			} else {
				if value == numberBefore {
					sameAdjacent = true
				}
			}
			numberBefore = value
		}
		return sameAdjacent
	}
}

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		log.Fatal("You must supply two numbers representing a range.")
	}
	left, _ := strconv.Atoi(args[0])
	right, _ := strconv.Atoi(args[1])
	var counter int = 0
	for candidate := left; candidate <= right; candidate++ {
		if meetCriteria(candidate) {
			counter++
		}
	}
	fmt.Printf("Passwords that meet criteria: %d\n", counter)
}
