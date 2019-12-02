// Álvaro Castellano Vela 2019/12/01
package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRunIntCodeExample1(t *testing.T) {

	input := []int{1, 0, 0, 0, 99}
	desiredOutput := []int{2, 0, 0, 0, 99}
	runIntCode(input)
	if reflect.DeepEqual(input, desiredOutput) == false {
		fmt.Printf("%v\n", input)
		fmt.Printf("%v\n", desiredOutput)
		t.Errorf("runIntCode not working as expected")

	}
}

func TestRunIntCodeExample2(t *testing.T) {

	input := []int{2, 3, 0, 3, 99}
	desiredOutput := []int{2, 3, 0, 6, 99}
	runIntCode(input)
	if reflect.DeepEqual(input, desiredOutput) == false {
		fmt.Printf("%v\n", input)
		fmt.Printf("%v\n", desiredOutput)
		t.Errorf("runIntCode not working as expected")
	}
}

func TestRunIntCodeExample3(t *testing.T) {

	input := []int{2, 4, 4, 5, 99, 0}
	desiredOutput := []int{2, 4, 4, 5, 99, 9801}
	runIntCode(input)
	if reflect.DeepEqual(input, desiredOutput) == false {
		fmt.Printf("%v\n", input)
		fmt.Printf("%v\n", desiredOutput)
		t.Errorf("runIntCode not working as expected")
	}
}

func TestRunIntCodeExample4(t *testing.T) {

	input := []int{1, 1, 1, 4, 99, 5, 6, 0, 99}
	desiredOutput := []int{30, 1, 1, 4, 2, 5, 6, 0, 99}
	runIntCode(input)
	if reflect.DeepEqual(input, desiredOutput) == false {
		fmt.Printf("%v\n", input)
		fmt.Printf("%v\n", desiredOutput)
		t.Errorf("runIntCode not working as expected")
	}
}
