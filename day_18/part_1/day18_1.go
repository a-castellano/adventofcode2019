// Álvaro Castellano Vela 2019/12/25
package main

import (
	"bufio"
	"fmt"
	"github.com/a-castellano/dijkstra"
	"log"
	//	"math/rand"
	"os"
	"regexp"
	//	"time"
	"strconv"
	//	"strings"
)

type Choice struct {
	vault        *dijkstra.Graph
	currentKeys  map[int]rune
	currentDoors map[int]rune
	distances    map[string]dijkstra.BestPath
	doorsInPaths map[int]map[string]bool
	initialPoint int
	formerStepts int
}

func processFile(filename string) ([][]rune, int, int, int) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var vault [][]rune
	var initialPoint int
	var i int = 0

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		var row []rune
		line := scanner.Text()
		for j, character := range line {
			row = append(row, character)
			if character == '@' {
				initialPoint = i*100 + j
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

func getKeysAndDoors(matrix [][]rune, rows int, columns int) (map[int]rune, map[int]rune) {

	keys := make(map[int]rune)
	doors := make(map[int]rune)

	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			if matrix[i][j] >= 'a' && matrix[i][j] <= 'z' {
				keys[i*100+j] = matrix[i][j]
			} else {
				if matrix[i][j] >= 'A' && matrix[i][j] <= 'Z' {
					doors[i*100+j] = matrix[i][j]
				}
			}
		}
	}

	return keys, doors
}

func generateGraph(matrix [][]rune, rows int, columns int) *dijkstra.Graph {
	graph := dijkstra.NewGraph()
	for i := 1; i < rows-1; i++ {
		for j := 1; j < columns-1; j++ {
			graph.AddVertex(i*100 + j)
		}
	}

	for i := 1; i < rows-1; i++ {
		for j := 1; j < columns-1; j++ {
			var formerCost int64 = 1
			if matrix[i][j] >= 'A' && matrix[i][j] <= 'Z' {
				formerCost += 100000
			}

			if matrix[i-1][j] != '#' {
				if matrix[i-1][j] <= 'A' && matrix[i-1][j] >= 'Z' {
					formerCost += 100000
				}
				graph.AddArc(i*100+j, (i-1)*100+j, formerCost)
			}
			if matrix[i+1][j] != '#' {
				if matrix[i+1][j] <= 'A' && matrix[i+1][j] >= 'Z' {
					formerCost += 100000
				}
				graph.AddArc(i*100+j, (i+1)*100+j, formerCost)
			}
			if matrix[i][j-1] != '#' {
				if matrix[i][j-1] <= 'A' && matrix[i][j-1] >= 'Z' {
					formerCost += 100000
				}
				graph.AddArc(i*100+j, i*100+j-1, formerCost)
			}
			if matrix[i][j+1] != '#' {
				if matrix[i][j+1] <= 'A' && matrix[i][j+1] >= 'Z' {
					formerCost += 100000
				}
				graph.AddArc(i*100+j, i*100+j+1, formerCost)
			}
		}
	}

	return graph
}

func findDoorInPath(bestPath dijkstra.BestPath, currentDoors map[int]rune) map[int]bool {
	doorsFound := make(map[int]bool)
	for _, point := range bestPath.Path {
		if _, ok := currentDoors[point]; ok {
			doorsFound[point] = true
		}
	}
	return doorsFound
}

func getPoints(pathID string) (int, int) {
	var origin, dest int
	re := regexp.MustCompile("([[:digit:]]+)-([[:digit:]]+)")
	inputString := re.FindAllStringSubmatch(pathID, -1)[0]
	origin, _ = strconv.Atoi(inputString[1])
	dest, _ = strconv.Atoi(inputString[2])
	return origin, dest
}

func getDistance(vault *dijkstra.Graph, origin int, dest int) dijkstra.BestPath {

	best, err := vault.Shortest(origin, dest)
	if err != nil {
		log.Fatal(err)
	}
	return best
}

func getDistances(vault *dijkstra.Graph, currentKeys map[int]rune, currentDoors map[int]rune, initialPoint int) (map[string]dijkstra.BestPath, map[int]map[string]bool) {
	distances := make(map[string]dijkstra.BestPath)
	doorsInPaths := make(map[int]map[string]bool)
	for point, _ := range currentKeys {
		pathID := fmt.Sprintf("%d-%d", initialPoint, point)
		best, err := vault.Shortest(initialPoint, point)
		if err != nil {
			log.Fatal(err)
		}
		distances[pathID] = best
		for door, _ := range findDoorInPath(best, currentDoors) {
			if _, ok := doorsInPaths[door]; !ok {
				doorsInPaths[door] = make(map[string]bool)
			}
			doorsInPaths[door][pathID] = true
		}
	}
	for origin, _ := range currentKeys {
		for destination, _ := range currentKeys {
			if origin != destination {

				pathID := fmt.Sprintf("%d-%d", origin, destination)
				best, err := vault.Shortest(origin, destination)
				if err != nil {
					log.Fatal(err)
				}

				distances[pathID] = best
				for door, _ := range findDoorInPath(best, currentDoors) {
					if _, ok := doorsInPaths[door]; !ok {
						doorsInPaths[door] = make(map[string]bool)
					}
					doorsInPaths[door][pathID] = true
				}

			}
		}
	}

	return distances, doorsInPaths
}

func copyVault(vault *dijkstra.Graph) *dijkstra.Graph {
	copied := dijkstra.NewGraph()

	for _, vertex := range vault.Verticies {
		if vertex.ID > 0 {
			copied.AddVertex(vertex.ID)
		}
	}

	for _, vertex := range vault.Verticies {
		arcs := vertex.GetArcs()
		if len(arcs) > 0 {
			for dest, distance := range arcs {
				copied.AddArc(vertex.ID, dest, distance)
			}
		}
	}
	return copied
}

func getReachableKeys(initialPoint int, currentKeys map[int]rune, distances map[string]dijkstra.BestPath) []int {

	var reachableKeys []int

	for point, _ := range currentKeys {
		pathID := fmt.Sprintf("%d-%d", initialPoint, point)
		if distance, ok := distances[pathID]; ok {
			if distance.Distance < 10000 {
				reachableKeys = append(reachableKeys, point)
			}
		}
	}
	return reachableKeys
}

func recalcuteArcs(vault *dijkstra.Graph, door int) {

	vertex := vault.Verticies[door]
	for point, value := range vertex.GetArcs() {
		vault.AddArc(door, point, value-100000)
	}

}

func copyMap(origin map[int]rune) map[int]rune {
	copied := make(map[int]rune)
	for index, element := range origin {
		copied[index] = element
	}
	return copied
}

func copyDistancesMap(origin map[string]dijkstra.BestPath) map[string]dijkstra.BestPath {
	copied := make(map[string]dijkstra.BestPath)
	for index, _ := range origin {
		var best dijkstra.BestPath
		best.Distance = origin[index].Distance
		best.Path = make([]int, len(origin[index].Path))
		for i, v := range origin[index].Path {
			best.Path[i] = v
		}
		copied[index] = best
	}
	return copied
}

func copyDoorsInPaths(origin map[int]map[string]bool) map[int]map[string]bool {
	copied := make(map[int]map[string]bool)
	for index, _ := range origin {
		copied[index] = make(map[string]bool)
		for path, _ := range origin[index] {
			copied[index][path] = true
		}
	}
	return copied
}

func generateChoices(numberOfChoices int, vault *dijkstra.Graph, currentKeys map[int]rune, currentDoors map[int]rune, distances map[string]dijkstra.BestPath, doorsInPaths map[int]map[string]bool, initialPoint int, formerStepts int) []Choice {
	var choices []Choice
	for i := 0; i < numberOfChoices; i++ {
		var newChoice Choice
		newChoice.vault = copyVault(vault)
		newChoice.currentKeys = copyMap(currentKeys)
		newChoice.currentDoors = copyMap(currentDoors)
		newChoice.distances = copyDistancesMap(distances)
		newChoice.doorsInPaths = copyDoorsInPaths(doorsInPaths)
		newChoice.initialPoint = initialPoint
		newChoice.formerStepts = formerStepts

		choices = append(choices, newChoice)
	}
	return choices
}

func walkPath(vault *dijkstra.Graph, currentKeys map[int]rune, currentDoors map[int]rune, distances map[string]dijkstra.BestPath, doorsInPaths map[int]map[string]bool, initialPoint int, keyPoint int, formerStepts int, bestPath int) (int, int, map[int]rune, map[int]rune, map[string]dijkstra.BestPath, map[int]map[string]bool, *dijkstra.Graph) {
	var recollectedKeys []int
	pathID := fmt.Sprintf("%d-%d", initialPoint, keyPoint)
	formerStepts += int(distances[pathID].Distance)
	if currentKeys[keyPoint] == 'b' {

	}
	for _, point := range distances[pathID].Path {
		if _, ok := currentKeys[point]; ok {
			recollectedKeys = append(recollectedKeys, point)
		}
	}
	// Unlock Doors
	for _, keyLocation := range recollectedKeys {
		var doorToUnlock int = -1
		for point, door := range currentDoors {
			if door == currentKeys[keyLocation]-32 {
				doorToUnlock = point
				break
			}
		}
		if doorToUnlock >= 0 {
			delete(currentDoors, doorToUnlock)
			// Recalculate Changed Distances
			// Distances affected by removed door
			recalcuteArcs(vault, doorToUnlock)
			for pathID, _ := range doorsInPaths[doorToUnlock] {
				//origin, dest := getPoints(pathID)
				fomerDistance := distances[pathID]
				fomerDistance.Distance -= 100000
				distances[pathID] = fomerDistance
				//distances[pathID] = getDistance(vault, origin, dest)

			}
			delete(doorsInPaths, doorToUnlock)
		}
		delete(currentKeys, keyLocation)
	}
	initialPoint = keyPoint
	return formerStepts, initialPoint, currentKeys, currentDoors, distances, doorsInPaths, vault
}

func getAllKeys(vault *dijkstra.Graph, currentKeys map[int]rune, currentDoors map[int]rune, distances map[string]dijkstra.BestPath, doorsInPaths map[int]map[string]bool, initialPoint int, formerStepts int, bestPath int) (int, int, *dijkstra.Graph, map[int]rune, map[int]rune, map[string]dijkstra.BestPath, map[int]map[string]bool, int) {
	if formerStepts < bestPath {
		reachableKeys := getReachableKeys(initialPoint, currentKeys, distances)
		if len(reachableKeys) == 0 && len(currentKeys) != 0 {
			fmt.Println("__________________________________________")
			fmt.Println("KEYS", currentKeys)
			fmt.Println("DOORS", currentDoors)
			fmt.Println("INITIALPOINT", initialPoint)
			log.Fatal("ABORT")
		}
		if len(reachableKeys) == 1 {
			keyPoint := reachableKeys[0]

			formerStepts, initialPoint, currentKeys, currentDoors, distances, doorsInPaths, vault = walkPath(vault, currentKeys, currentDoors, distances, doorsInPaths, initialPoint, keyPoint, formerStepts, bestPath)
		} else {
			if len(reachableKeys) > 1 {
				reachableKeysLength := len(reachableKeys)
				//if reachableKeysLength > 3 {
				//	reachableKeysLength = 3
				//}
				choices := generateChoices(reachableKeysLength, vault, currentKeys, currentDoors, distances, doorsInPaths, initialPoint, formerStepts)
				for i := 0; i < len(choices); i++ {

					keyPoint := reachableKeys[i]

					choices[i].formerStepts, choices[i].initialPoint, choices[i].currentKeys, choices[i].currentDoors, choices[i].distances, choices[i].doorsInPaths, choices[i].vault = walkPath(choices[i].vault, choices[i].currentKeys, choices[i].currentDoors, choices[i].distances, choices[i].doorsInPaths, choices[i].initialPoint, keyPoint, choices[i].formerStepts, bestPath)

					choiceReachableKeys := getReachableKeys(choices[i].initialPoint, choices[i].currentKeys, choices[i].distances)
					if len(choiceReachableKeys) > 0 {
						choices[i].formerStepts, choices[i].initialPoint, choices[i].vault, choices[i].currentKeys, choices[i].currentDoors, choices[i].distances, choices[i].doorsInPaths, _ = getAllKeys(choices[i].vault, choices[i].currentKeys, choices[i].currentDoors, choices[i].distances, choices[i].doorsInPaths, choices[i].initialPoint, choices[i].formerStepts, bestPath)
					}
				}
				bestChoice := -1
				bestFormerStepts := bestPath
				for i := 0; i < len(choices); i++ {
					if choices[i].formerStepts < bestFormerStepts {
						bestFormerStepts = choices[i].formerStepts
						bestChoice = i
					}
				}
				if bestChoice == -1 {
					bestChoice = 0
					choices[bestChoice].formerStepts = 10000000000
				}
				vault = choices[bestChoice].vault
				currentKeys = choices[bestChoice].currentKeys
				currentDoors = choices[bestChoice].currentDoors
				distances = choices[bestChoice].distances
				doorsInPaths = choices[bestChoice].doorsInPaths
				initialPoint = choices[bestChoice].initialPoint
				formerStepts = choices[bestChoice].formerStepts
			}
		}
		if len(currentKeys) > 0 {
			formerStepts, initialPoint, vault, currentKeys, currentDoors, distances, doorsInPaths, bestPath = getAllKeys(vault, currentKeys, currentDoors, distances, doorsInPaths, initialPoint, formerStepts, bestPath)
		} else {
			if formerStepts >= bestPath {
				formerStepts = 10000000000
				return formerStepts, initialPoint, vault, currentKeys, currentDoors, distances, doorsInPaths, bestPath
			} else {
				bestPath = formerStepts
			}
		}
	}

	return formerStepts, initialPoint, vault, currentKeys, currentDoors, distances, doorsInPaths, bestPath
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	matrix, rows, columns, initialPoint := processFile(filename)
	vault := generateGraph(matrix, rows, columns)
	currentKeys, currentDoors := getKeysAndDoors(matrix, rows, columns)
	distances, doorsInPaths := getDistances(vault, currentKeys, currentDoors, initialPoint)
	stepts, _, _, _, _, _, _, _ := getAllKeys(vault, currentKeys, currentDoors, distances, doorsInPaths, initialPoint, 0, 140)
	fmt.Println("Steps required: ", stepts)
}
