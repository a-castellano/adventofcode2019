// Álvaro Castellano Vela 2019/12/21
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	//	"reflect"
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

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func processFile(filename string) ([]int, int) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()

	line := scanner.Text()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile("[[:digit:]]")
	inputString := re.FindAllStringSubmatch(line, -1)
	inputStringSize := len(inputString)

	var input []int

	for i := 0; i < 10000; i++ {
		for _, str := range inputString {
			digit, _ := strconv.Atoi(str[0])
			input = append(input, digit)
		}
	}

	return input, inputStringSize * 10000
}

func calculatePhase(input []int, inputLenght int, initialOffset int) []int {

	newInput := make([]int, inputLenght)

	newInput[inputLenght-1] = input[inputLenght-1]
	for i := inputLenght - 2; i >= initialOffset; i-- {

		newInput[i] = (newInput[i+1] + input[i]) % 10
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
	input, inputLenght := processFile(filename)
	var initialOffset int = input[0]*1000000 + input[1]*100000 + input[2]*10000 + input[3]*1000 + input[4]*100 + input[5]*10 + input[6]
	for i := 0; i < phases; i++ {
		input = calculatePhase(input, inputLenght, initialOffset)
	}
	fmt.Println(input[initialOffset : initialOffset+8])
}
