// Álvaro Castellano Vela 2019/12/25
package main

import (
	"bufio"
	"fmt"
	"github.com/RyanCarrier/dijkstra"
	"log"
	//	"math/rand"
	"os"
	//	"time"
	//	"strconv"
	//	"strings"
)

type Point struct {
	i int
	j int
}

type Choice struct {
	vault        [][]rune
	currentKeys  map[rune]Point
	currentDoors map[rune]Point
	initialPoint Point
	formerStepts int
}

func processFile(filename string) ([][]rune, int, int, Point) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var vault [][]rune
	var initialPoint Point
	var i int = 0

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		var row []rune
		line := scanner.Text()
		for j, character := range line {
			row = append(row, character)
			if character == '@' {
				initialPoint.i = i
				initialPoint.j = j
			}
		}
		vault = append(vault, row)
		i++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return vault, len(vault), len(vault[0]), initialPoint
}

func getKeysAndDoors(vault [][]rune, rows int, columns int) (map[rune]Point, map[rune]Point) {

	keys := make(map[rune]Point)
	doors := make(map[rune]Point)

	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			if vault[i][j] >= 'a' && vault[i][j] <= 'z' {
				newPoint := Point{i, j}
				keys[vault[i][j]] = newPoint
			} else {
				if vault[i][j] >= 'A' && vault[i][j] <= 'Z' {
					newPoint := Point{i, j}
					doors[vault[i][j]] = newPoint
				}
			}
		}
	}

	return keys, doors
}

func generateGraph(vault [][]rune, rows int, columns int) *dijkstra.Graph {
	graph := dijkstra.NewGraph()
	for i := 1; i < rows-1; i++ {
		for j := 1; j < columns-1; j++ {
			graph.AddVertex(i*1000 + j)
		}
	}

	for i := 1; i < rows-1; i++ {
		for j := 1; j < columns-1; j++ {
			var formerCost int64 = 1
			if vault[i][j] >= 'A' && vault[i][j] <= 'Z' {
				formerCost += 100000
			}

			if vault[i-1][j] != '#' {
				if vault[i-1][j] <= 'A' && vault[i-1][j] >= 'Z' {
					formerCost += 100000
				}
				graph.AddArc(i*1000+j, (i-1)*1000+j, formerCost)
			}
			if vault[i+1][j] != '#' {
				if vault[i+1][j] <= 'A' && vault[i+1][j] >= 'Z' {
					formerCost += 100000
				}
				graph.AddArc(i*1000+j, (i+1)*1000+j, formerCost)
			}
			if vault[i][j-1] != '#' {
				if vault[i][j-1] <= 'A' && vault[i][j-1] >= 'Z' {
					formerCost += 100000
				}
				graph.AddArc(i*1000+j, i*1000+j-1, formerCost)
			}
			if vault[i][j+1] != '#' {
				if vault[i][j+1] <= 'A' && vault[i][j+1] >= 'Z' {
					formerCost += 100000
				}
				graph.AddArc(i*1000+j, i*1000+j+1, formerCost)
			}
		}
	}

	return graph
}

func getPointFromID(id int) Point {
	var point Point
	point.i = id / 1000
	point.j = id % 1000

	return point
}

func walkPath(vault [][]rune, currentKeys map[rune]Point, currentDoors map[rune]Point, path []int, initialPoint Point) (Point, [][]rune) {
	//fmt.Println(path)
	vault[initialPoint.i][initialPoint.j] = '.'
	for _, step := range path {
		point := getPointFromID(step)
		//	fmt.Println(point)
		if vault[point.i][point.j] != '.' {
			//		fmt.Printf("key found: %c\n", vault[point.i][point.j])
			delete(currentKeys, vault[point.i][point.j])
			//		fmt.Println(currentDoors)
			//		fmt.Println(vault[point.i][point.j])
			//		fmt.Println(currentDoors[vault[point.i][point.j]-32])
			doorPoint := currentDoors[vault[point.i][point.j]-32]
			vault[doorPoint.i][doorPoint.j] = '.'
			delete(currentDoors, vault[point.i][point.j]-32)
			delete(currentKeys, vault[point.i][point.j])
			initialPoint.i = point.i
			initialPoint.j = point.j
			vault[point.i][point.j] = '.'
		}
	}
	vault[initialPoint.i][initialPoint.j] = '@'
	return initialPoint, vault
}

func printVault(vault [][]rune, rows int, columns int) {

	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			fmt.Printf("%c", vault[i][j])
		}
		fmt.Printf("\n")
	}

}

func copyMap(origin map[rune]Point) map[rune]Point {
	copied := make(map[rune]Point)
	for index, element := range origin {
		copied[index] = element
	}
	return copied
}

func copyVault(origin [][]rune, rows int, columns int) [][]rune {
	copied := make([][]rune, rows)
	for i := range origin {
		copied[i] = make([]rune, columns)
		copy(copied[i], origin[i])
	}
	return copied
}

func generateChoices(vault [][]rune, rows int, columns int, numberOfChoices int, currentKeys map[rune]Point, currentDoors map[rune]Point, initialPoint Point, formerStepts int) []Choice {
	var choices []Choice
	for i := 0; i < numberOfChoices; i++ {
		var newChoice Choice
		newChoice.currentDoors = copyMap(currentDoors)
		newChoice.currentKeys = copyMap(currentKeys)
		newChoice.vault = copyVault(vault, rows, columns)
		newChoice.initialPoint = initialPoint
		newChoice.formerStepts = formerStepts
		choices = append(choices, newChoice)
	}
	return choices
}

func getKeys(vault [][]rune, rows int, columns int, initialPoint Point, formerStepts int, bestChoiceValue int, cache map[string]dijkstra.BestPath) ([][]rune, int, int, Point, int) {

	//fmt.Println(":::::::::::")
	//printVault(vault, rows, columns)
	//fmt.Println(":::::::::::")
	if formerStepts < bestChoiceValue {

		currentKeys, currentDoors := getKeysAndDoors(vault, rows, columns)
		graph := generateGraph(vault, rows, columns)

		reachableKeys := make(map[rune]dijkstra.BestPath)

		//fmt.Println(initialPoint)
		//fmt.Println(currentKeys)
		//fmt.Println(currentDoors)
		initialPointID := initialPoint.i*1000 + initialPoint.j
		//Calculate reachable keys
		//s1 := rand.NewSource(time.Now().UnixNano())
		//r1 := rand.New(s1)
		//randID := r1.Intn(100)
		for key, point := range currentKeys {
			pointID := point.i*1000 + point.j
			pathID := fmt.Sprintf("%d-%d", initialPointID, pointID)
			fmt.Println("PATHID ", pathID)
			//	best, err := graph.Shortest(initialPointID, pointID)
			//	if err != nil {
			//		log.Fatal(err)
			//	}
			//	if best.Distance < 10000 {
			//		fmt.Println(randID, "Shortest distance from ", initialPointID, " to ", string(key), "in Point", pointID, " : ", best.Distance, " following path ", best.Path)
			//		reachableKeys[key] = best
			//reachableKeys[key] = cache[pathID]
			//			}
			if _, ok := cache[pathID]; ok {
				reachableKeys[key] = cache[pathID]
				fmt.Println("Value for ", pathID, " is ", cache[pathID].Distance)
				fmt.Println("Value for ", pathID, " is ", cache[pathID].Path)
				//				best, err := graph.Shortest(initialPointID, pointID)
				//				if err != nil {
				//					log.Fatal(err)
				//				}
				//				if best.Distance < 10000 {
				//					//	fmt.Println("Shortest distance from initialPoint to ", string(key), " : ", best.Distance, " following path ", best.Path)
				//					fmt.Println("BUT ", pathID, best.Distance)
				//					reachableKeys[key] = best
				//					if best.Distance != cache[pathID].Distance {
				//						log.Fatal("dddd")
				//					}
				//				}
			} else {
				//	fmt.Println("Looking for pathID: ", pathID)
				best, err := graph.Shortest(initialPointID, pointID)
				if err != nil {
					log.Fatal(err)
				}
				if best.Distance < 10000 {
					cache[pathID] = best
					reachableKeys[key] = cache[pathID]
				}

			}
		}

		keys := make([]rune, len(reachableKeys))
		i := 0
		for key, _ := range reachableKeys {
			keys[i] = key
			i++
		}
		//fmt.Println("KEYS:", reachableKeys)
		if len(reachableKeys) == 1 {
			//fmt.Println("There is only one key to reach")
			//fmt.Println(reachableKeys[keys[0]])
			initialPoint, vault = walkPath(vault, currentKeys, currentDoors, reachableKeys[keys[0]].Path, initialPoint)
			fmt.Println(formerStepts, "__________", reachableKeys[keys[0]].Distance)
			formerStepts += int(reachableKeys[keys[0]].Distance)
			fmt.Println(formerStepts, "__________")

			printVault(vault, rows, columns)
			currentKeys, currentDoors = getKeysAndDoors(vault, rows, columns)
		} else {
			if len(reachableKeys) > 1 {
				//		fmt.Println("Current steps: ", formerStepts)
				//		fmt.Println("::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::")
				//		printVault(vault, rows, columns)
				choices := generateChoices(vault, rows, columns, len(reachableKeys), currentKeys, currentDoors, initialPoint, formerStepts)
				var goodLookingChoices []int
				//fmt.Println()
				for i := 0; i < len(keys); i++ {
					//	fmt.Println(choices[i])
					//	fmt.Println()

					choices[i].initialPoint, choices[i].vault = walkPath(choices[i].vault, choices[i].currentKeys, choices[i].currentDoors, reachableKeys[keys[i]].Path, choices[i].initialPoint)
					choices[i].formerStepts += int(reachableKeys[keys[i]].Distance)
					choices[i].vault, rows, columns, choices[i].initialPoint, choices[i].formerStepts = getKeys(choices[i].vault, rows, columns, choices[i].initialPoint, choices[i].formerStepts, bestChoiceValue, cache)
					if choices[i].formerStepts <= bestChoiceValue {
						bestChoiceValue = choices[i].formerStepts
						goodLookingChoices = append(goodLookingChoices, i)
					}
					choices[i].currentKeys, choices[i].currentDoors = getKeysAndDoors(choices[i].vault, rows, columns)
					if formerStepts == 0 {
						fmt.Println("Choice ", i, " total steps", choices[i].formerStepts)
					}
				}

				var minimunDistance int = 10000000000
				var bestChoice int = -1
				for i := 0; i < len(goodLookingChoices); i++ {
					//fmt.Println("Choice: ", i, " steps: ", choices[i].formerStepts)
					if choices[goodLookingChoices[i]].formerStepts < minimunDistance {
						minimunDistance = choices[goodLookingChoices[i]].formerStepts
						bestChoice = goodLookingChoices[i]
					}
				}
				fmt.Println("__________CHECKING CHOICES, there are ", len(keys), " choices. Best one is", bestChoice)
				if bestChoice != -1 {
					vault = choices[bestChoice].vault
					initialPoint = choices[bestChoice].initialPoint
					currentKeys = choices[bestChoice].currentKeys
					currentDoors = choices[bestChoice].currentDoors
					formerStepts = choices[bestChoice].formerStepts
					//bestChoiceValue = formerStepts
					currentKeys, currentDoors = getKeysAndDoors(vault, rows, columns)
					fmt.Println("Choice ", bestChoice, "of ", len(goodLookingChoices), "-  total steps", choices[bestChoice].formerStepts)
				} else {
					vault = choices[0].vault
					initialPoint = choices[0].initialPoint
					currentKeys = choices[0].currentKeys
					currentDoors = choices[0].currentDoors
					formerStepts = 1000000000000
					currentKeys, currentDoors = getKeysAndDoors(vault, rows, columns)

				}
				//fmt.Println("formerStepts: ", formerStepts)
				//fmt.Println("currentKeys: ", currentKeys)

				//		fmt.Println("::::::::::::::::::::::::::::::::::::::::::::::::")
				//		printVault(choices[bestChoice].vault, rows, columns)
				//		printVault(vault, rows, columns)
				//		fmt.Println("::::::::::::::::::::::::::::::::::::::::::::::::")

				//		fmt.Println("_______________________________________________________")
				//		fmt.Println("_______________________________________________________")
				//		fmt.Println("_______________________________________________________")
			}
		}

		currentKeys, currentDoors = getKeysAndDoors(vault, rows, columns)
		if len(currentKeys) > 0 {

			if formerStepts <= bestChoiceValue {
				vault, rows, columns, initialPoint, formerStepts = getKeys(vault, rows, columns, initialPoint, formerStepts, bestChoiceValue, cache)
			} else {
				return vault, rows, columns, initialPoint, 10000000000000000
			}
		} // else {
		//   fmt.Println("$$$$$$$$$$$$$$$")
		//   printVault(vault, rows, columns)
		//   fmt.Println("$$$$$$$$$$$$$$$444")

		//   fmt.Println("Path finished with ", formerStepts, " steps.")
		//   fmt.Println("Best choice is  ", bestChoiceValue)
		//}//
	}
	return vault, rows, columns, initialPoint, formerStepts
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	vault, rows, columns, initialPoint := processFile(filename)
	cache := make(map[string]dijkstra.BestPath)
	//fmt.Println("Result: ", vault)
	//fmt.Println("Rows: ", rows)
	//fmt.Println("Columns: ", columns)
	//fmt.Println("initialPoint: ", initialPoint)
	_, _, _, _, stepts := getKeys(vault, rows, columns, initialPoint, 0, 1000000, cache)
	fmt.Println("Shortest path: ", stepts)
}
