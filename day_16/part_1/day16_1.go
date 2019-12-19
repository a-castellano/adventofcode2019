// Álvaro Castellano Vela 2019/12/19
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

func processFile(filename string) []int {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var input []int

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()

	line := scanner.Text()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile("[[:digit:]]")
	inputString := re.FindAllStringSubmatch(line, -1)
	for _, str := range inputString {
		digit, _ := strconv.Atoi(str[0])
		input = append(input, digit)
	}

	return input
}

func createPaterns(input []int) ([][]int, int) {
	basePatern := [...]int{0, 1, 0, -1}
	var patterns [][]int

	var inputLenght int = len(input)
	for i := 1; i <= inputLenght; i++ {
		var pattern []int
		for _, patternDigit := range basePatern {
			for j := 0; j < i; j++ {
				pattern = append(pattern, patternDigit)
			}
		}
		patterns = append(patterns, pattern)
	}
	return patterns, inputLenght
}

func calculatePhase(input []int, inputLenght int, patterns [][]int) []int {

	var newInput []int

	for i := 0; i < inputLenght; i++ {

		var inputDigit int
		for j := 0; j < inputLenght; j++ {

			var patternLenght int = len(patterns[i])
			inputDigit += input[j] * patterns[i][(1+j)%patternLenght]
		}
		inputDigit = Abs(inputDigit)
		for inputDigit >= 10 {
			inputDigit = inputDigit % 10
		}
		newInput = append(newInput, inputDigit)
	}
	return newInput
}

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		log.Fatal("You must supply a file to process and the number of phases you want to calculate.")
	}
	filename := args[0]
	phases, _ := strconv.Atoi(args[1])
	input := processFile(filename)
	patterns, inputLenght := createPaterns(input)
	for i := 0; i < phases; i++ {
		input = calculatePhase(input, inputLenght, patterns)
	}
	fmt.Println(input[:8])
}
