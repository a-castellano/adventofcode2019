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
		quantityGenerated, outpuChemicalName := getChemicalInfo(outputChemicalString)
		if _, ok := chemicals[outpuChemicalName]; !ok {
			outputChemical := Chemical{outpuChemicalName, quantityGenerated, make([]RequiredChemical, 0)}
			chemicals[outpuChemicalName] = &outputChemical
		} else {
			chemicals[outpuChemicalName].quantityGenerated = quantityGenerated
		}
		//fmt.Println(quantityGenerated, outpuChemicalName)
		for _, inputString := range inputChemicalsString {
			quantityRequired, inputChemicalName := getChemicalInfo(inputString[0])
			if _, ok := chemicals[inputChemicalName]; !ok {
				inputChemical := Chemical{inputChemicalName, -1, make([]RequiredChemical, 0)}
				chemicals[inputChemicalName] = &inputChemical
			}
			requiredChemical := RequiredChemical{chemicals[inputChemicalName], quantityRequired}
			chemicals[outpuChemicalName].requiredChemicals = append(chemicals[outpuChemicalName].requiredChemicals, requiredChemical)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return chemicals
}

func generateChemical(chemicals map[string]*Chemical, requiredChemicalName string, requiredQuantity int, state State, decision map[string]int) State {

	outputChemical := chemicals[requiredChemicalName]
	//Check if this chemicals only needs ORE to be generated
	var requiredChemicalsLenght int = len(outputChemical.requiredChemicals)
	if requiredChemicalsLenght == 1 && outputChemical.requiredChemicals[0].chemical.name == "ORE" {
		var requiredChemicals int = outputChemical.requiredChemicals[0].quantity
		times := 0
		if state.boxOfChemicals["ORE"] < requiredChemicals {
			generated := 0
			for generated < requiredQuantity {
				state.generatedChemicals["ORE"] += requiredChemicals
				state.boxOfChemicals["ORE"] += requiredChemicals
				state.generatedChemicals[outputChemical.name] += outputChemical.quantityGenerated
				generated += outputChemical.quantityGenerated
				state.boxOfChemicals[outputChemical.name] += outputChemical.quantityGenerated
				times++
			}
		}
		state.boxOfChemicals["ORE"] -= requiredChemicals * times
	} else {
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
			operation := fmt.Sprintf("%v", outputChemical.requiredChemicals)

			permutations := permutation(rangeSlice(0, len(outputChemical.requiredChemicals)))

			if _, ok := decision[operation]; ok {
				requiredPermutation := permutations[decision[operation]]
				var optimizedPermutations [][]int
				optimizedPermutations = append(optimizedPermutations, requiredPermutation)
				permutations = optimizedPermutations
			}
			var posibleStates []State
			for i := 0; i < len(permutations); i++ {
				posibleStates = append(posibleStates, copyState(state))
			}
			var realCandidates []State
			bestPath := 999999999999
			bestPremutationIndex := -1
			for premutationIndex, path := range permutations {
				currentState := posibleStates[premutationIndex]
				for _, index := range path {
					requiredChemical := outputChemical.requiredChemicals[index]
					if currentState.boxOfChemicals[requiredChemical.chemical.name] >= requiredChemical.quantity {
						currentState.boxOfChemicals[requiredChemical.chemical.name] -= requiredChemical.quantity
					} else {

						for currentState.boxOfChemicals[requiredChemical.chemical.name] < requiredChemical.quantity && currentState.generatedChemicals["ORE"] < bestPath {
							var chemicalRequiredQuantity int = requiredChemical.quantity - currentState.boxOfChemicals[requiredChemical.chemical.name]
							currentState = generateChemical(chemicals, requiredChemical.chemical.name, chemicalRequiredQuantity, currentState, decision)
						}
						currentState.boxOfChemicals[requiredChemical.chemical.name] -= requiredChemical.quantity
						if currentState.generatedChemicals["ORE"] >= bestPath {
							currentState.generatedChemicals["ORE"] = 999999999999
							break
						}
					}
					if currentState.boxOfChemicals[requiredChemical.chemical.name] < 0 {
						currentState.generatedChemicals["ORE"] = 999999999999
						break
					}

				}
				if bestPath > currentState.generatedChemicals["ORE"] {
					bestPath = currentState.generatedChemicals["ORE"]
					realCandidates = append(realCandidates, currentState)
					bestPremutationIndex = premutationIndex
				}
			}

			if _, ok := decision[operation]; !ok {
				decision[operation] = bestPremutationIndex
			}
			state = realCandidates[0]
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

	decision := make(map[string]int)

	for _, chemical := range chemicals {
		boxOfChemicals[chemical.name] = 0
		generatedChemicals[chemical.name] = 0
	}

	state.boxOfChemicals = boxOfChemicals
	state.generatedChemicals = generatedChemicals
	state.minimunOreRequired = 999999999999

	state = generateChemical(chemicals, "FUEL", 1, state, decision)

	fmt.Println("boxOfChemicals", state.boxOfChemicals)
	fmt.Println("generatedChemicals", state.generatedChemicals)
}
