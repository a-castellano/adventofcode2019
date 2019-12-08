// Álvaro Castellano Vela 2019/12/08
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func processFile(filename string) []int {
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

	stringPixels := strings.Split(line, "")
	var pixels []int

	for _, value := range stringPixels {
		data, _ := strconv.Atoi(value)
		pixels = append(pixels, data)
	}

	return pixels
}

func findLayerWithFewestZeros(pixels []int, pixels_wide int, pixels_tall int) int {

	var pixel_size int = pixels_wide * pixels_tall

	var result int = 0
	var lewestZeros int = pixel_size + 1

	for pos := 0; pos < len(pixels); pos += pixel_size {
		var zeros, ones, twos int = 0, 0, 0
		for posNumber := 0; posNumber < pixel_size; posNumber++ {
			switch pixels[pos+posNumber] {
			case 0:
				zeros++
			case 1:
				ones++
			case 2:
				twos++
			}
		}
		if zeros < lewestZeros {
			lewestZeros = zeros
			result = ones * twos
		}
	}

	return result
}

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		log.Fatal("You must supply a file to process and Image size (pixels wide and pixels tall)")
	}
	filename := args[0]
	pixels_wide, _ := strconv.Atoi(args[1])
	pixels_tall, _ := strconv.Atoi(args[2])
	pixels := processFile(filename)
	result := findLayerWithFewestZeros(pixels, pixels_wide, pixels_tall)
	fmt.Printf("Result: %d\n", result)
}
