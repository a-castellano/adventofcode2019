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

func generateImage(pixels []int, pixels_wide int, pixels_tall int) {

	var pixel_size int = pixels_wide * pixels_tall

	imageInt := make([][]int, pixels_tall)
	for i := 0; i < pixels_tall; i++ {
		imageInt[i] = make([]int, pixels_wide)
		for j := 0; j < pixels_wide; j++ {
			imageInt[i][j] = 3
		}
	}

	image := make([][]rune, pixels_tall)
	for i := 0; i < pixels_tall; i++ {
		image[i] = make([]rune, pixels_wide)
	}

	for pos := len(pixels) - pixel_size; pos >= 0; pos -= pixel_size {

		row_counter := 0
		for posNumber := pos; posNumber < pos+pixel_size; posNumber += pixels_wide {
			for columnPos := posNumber; columnPos < posNumber+pixels_wide; columnPos++ {
				row := posNumber % pixels_tall
				column := columnPos % pixels_wide
				value := pixels[columnPos]
				switch imageInt[row][column] {
				case 0: // Black
					if value != 2 { // Transparent
						imageInt[row][column] = value
					}
				case 1: // White
					if value != 2 { // Transparent
						imageInt[row][column] = value
					}
				default: // Transparent
					imageInt[row][column] = value
				}
			}
			row_counter++
		}
	}
	for i := 0; i < pixels_tall; i++ {
		for j := 0; j < pixels_wide; j++ {
			if imageInt[i][j] == 0 {
				image[i][j] = ' '
			} else {
				image[i][j] = '#'
			}

			fmt.Printf("%c ", image[i][j])
		}
		fmt.Printf("\n")
	}
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
	generateImage(pixels, pixels_wide, pixels_tall)
}
