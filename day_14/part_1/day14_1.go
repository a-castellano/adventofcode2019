// Álvaro Castellano Vela 2019/12/14
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func copyState(target State) State {
	copiedBox := make(map[string]int)
	copiedGenerated := make(map[string]int)

	for key, value := range target.boxOfChemicals {
		copiedBox[key] = value
	}

	for key, value := range target.generatedChemicals {
		copiedGenerated[key] = value
	}

	return State{copiedBox, copiedGenerated, target.minimunOreRequired}
}

type State struct {
	boxOfChemicals     map[string]int
	generatedChemicals map[string]int
	minimunOreRequired int
}

func rangeSlice(start, stop int) []int {
	xs := make([]int, stop-start)
	for i := 0; i < len(xs); i++ {
		xs[i] = i + start
	}
	return xs
}

func permutation(xs []int) (permuts [][]int) {
	var rc func([]int, int)
	rc = func(a []int, k int) {
		if k == len(a) {
			permuts = append(permuts, append([]int{}, a...))
		} else {
			for i := k; i < len(xs); i++ {
				a[k], a[i] = a[i], a[k]
				rc(a, k+1)
				a[k], a[i] = a[i], a[k]
			}
		}
	}
	rc(xs, 0)

	return permuts
}

type Chemical struct {
	name              string
	quantityGenerated int
	requiredChemicals []RequiredChemical
}

type RequiredChemical struct {
	chemical *Chemical
	quantity int
}

//func (moon *Moon) calculatePotentialEnergy() {
//	moon.potentialEnergy = Abs(moon.position.X) + Abs(moon.position.Y) + Abs(moon.position.Z)
//}

func getChemicalInfo(info string) (int, string) {
	managedInfo := strings.Split(info, " ")
	quantity, _ := strconv.Atoi(managedInfo[0])
	chemicalName := managedInfo[1]

	return quantity, chemicalName
}

func processFile(filename string) map[string]*Chemical {

	chemicals := make(map[string]*Chemical)

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		re := regexp.MustCompile("[[:digit:]]+ [[:upper:]]+")
		match := re.FindAllStringSubmatch(scanner.Text(), -1)
		matchLength := len(match)
		outputChemicalString := match[matchLength-1][0]
		inputChemicalsString := match[:matchLength-1]
		fmt.Println("line: ", match)
		fmt.Println("output: ", outputChemicalString)
		fmt.Println("input: ", inputChemicalsString)
		quantityGenerated, outpuChemicalName := getChemicalInfo(outputChemicalString)
		if _, ok := chemicals[outpuChemicalName]; !ok {
			outputChemical := Chemical{outpuChemicalName, quantityGenerated, make([]RequiredChemical, 0)}
			chemicals[outpuChemicalName] = &outputChemical
		} else {
			chemicals[outpuChemicalName].quantityGenerated = quantityGenerated
		}
		fmt.Println(quantityGenerated, outpuChemicalName)
		for _, inputString := range inputChemicalsString {
			fmt.Println("Requires: ", inputString[0])
			quantityRequired, inputChemicalName := getChemicalInfo(inputString[0])
			fmt.Println("Requires: ", quantityRequired, inputChemicalName)
			if _, ok := chemicals[inputChemicalName]; !ok {
				fmt.Println(inputChemicalName, " not registered yet")
				inputChemical := Chemical{inputChemicalName, -1, make([]RequiredChemical, 0)}
				fmt.Println(inputChemical)
				chemicals[inputChemicalName] = &inputChemical
			}
			requiredChemical := RequiredChemical{chemicals[inputChemicalName], quantityRequired}
			chemicals[outpuChemicalName].requiredChemicals = append(chemicals[outpuChemicalName].requiredChemicals, requiredChemical)
		}
		fmt.Println("_________")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(chemicals)

	for key, _ := range chemicals {
		fmt.Println("Output: ", key)
		fmt.Println("Requires:")
		for _, required := range chemicals[key].requiredChemicals {
			fmt.Println(required.quantity, required.chemical.name)
		}
		fmt.Println("To generate ", chemicals[key].quantityGenerated)
		fmt.Println("")
	}
	return chemicals
}

func generateChemical(chemicals map[string]*Chemical, requiredChemicalName string, requiredQuantity int, state State) State {

	outputChemical := chemicals[requiredChemicalName]
	for _, requiredChemicals := range outputChemical.requiredChemicals {
		fmt.Println("Generate ", requiredQuantity, " ", requiredChemicalName, "requires", requiredChemicals.quantity, " of ", requiredChemicals.chemical.name)
	}
	//Check if this chemicals only needs ORE to be generated
	var requiredChemicalsLenght int = len(outputChemical.requiredChemicals)
	if requiredChemicalsLenght == 1 && outputChemical.requiredChemicals[0].chemical.name == "ORE" {
		var requiredChemicals int = outputChemical.requiredChemicals[0].quantity
		fmt.Println("ORE FOUND")
		fmt.Println(requiredChemicalName, " requires ", requiredChemicals, "of ORE")
		if state.boxOfChemicals["ORE"] < requiredChemicals {
			// More ORE is needed
			fmt.Println("Not enough ORE, generating more.")
			//generate more
			for state.boxOfChemicals["ORE"] < requiredChemicals {
				state.generatedChemicals["ORE"] += requiredChemicals
				state.boxOfChemicals["ORE"] += requiredChemicals
				state.generatedChemicals[outputChemical.name] += outputChemical.quantityGenerated
				state.boxOfChemicals[outputChemical.name] += outputChemical.quantityGenerated
			}
		}
		state.boxOfChemicals["ORE"] -= requiredChemicals
	} else {
		fmt.Println("NO NEED ORE DIRECTLY")
		// check if we need to generate more chemicals
		var generateMoreChemicals bool = false
		for _, requiredChemicals := range outputChemical.requiredChemicals {
			if state.boxOfChemicals[requiredChemicals.chemical.name] < requiredChemicals.quantity {
				generateMoreChemicals = true
				break
			}
		}
		if generateMoreChemicals == false {
			for _, requiredChemical := range outputChemical.requiredChemicals {
				// substract from box
				var requiredChemicals int = requiredChemical.quantity
				state.boxOfChemicals[requiredChemical.chemical.name] -= requiredChemicals
			}
		} else {
			fmt.Println("FUN")

			permutations := permutation(rangeSlice(0, len(outputChemical.requiredChemicals)))
			fmt.Println(permutations)
			var posibleStates []State
			for i := 0; i < len(permutations); i++ {
				// Generate a copy of boxOfChemicals and generatedChemicals
				posibleStates = append(posibleStates, copyState(state))
			}
			for premutationIndex, path := range permutations {
				currentState := posibleStates[premutationIndex]
				fmt.Println(path)
				for _, index := range path {
					// Check if we have enough of this chemical
					requiredChemical := outputChemical.requiredChemicals[index]
					fmt.Println("Check if we need ", requiredChemical.chemical.name)
					if currentState.boxOfChemicals[requiredChemical.chemical.name] < requiredChemical.quantity {
						fmt.Println("We need more ", requiredChemical.chemical.name, " chemical")
						fmt.Println("current state: ", currentState)
						var requiredQuantity int = requiredChemical.quantity - currentState.boxOfChemicals[requiredChemical.chemical.name]
						currentState = generateChemical(chemicals, requiredChemical.chemical.name, requiredQuantity, currentState)
						fmt.Println("current state after calling generateChemical: ", currentState)
					}
				}
			}
			//Choose best path
			bestPathValue := 999999999999999999
			bestPathIndex := -1
			for premutationIndex, _ := range permutations {
				if posibleStates[premutationIndex].boxOfChemicals["ORE"] < bestPathValue {
					bestPathValue = posibleStates[premutationIndex].boxOfChemicals["ORE"]
					bestPathIndex = premutationIndex
				}
			}
			state = posibleStates[bestPathIndex]
			fmt.Println("interfinal state: ", state)
		}
		state.boxOfChemicals[outputChemical.name] += outputChemical.quantityGenerated
	}
	return state
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]

	chemicals := processFile(filename)
	var state State
	boxOfChemicals := make(map[string]int)
	generatedChemicals := make(map[string]int)

	for _, chemical := range chemicals {
		boxOfChemicals[chemical.name] = 0
		generatedChemicals[chemical.name] = 0
	}

	state.boxOfChemicals = boxOfChemicals
	state.generatedChemicals = generatedChemicals
	state.minimunOreRequired = 999999999999999999

	fmt.Println("=====================================================")
	state = generateChemical(chemicals, "C", 1, state)

	fmt.Println("=====================================================")
	fmt.Println("=====================================================")
	fmt.Println("=====================================================")
	fmt.Println("boxOfChemicals", state.boxOfChemicals)
	fmt.Println("generatedChemicals", state.generatedChemicals)
}
