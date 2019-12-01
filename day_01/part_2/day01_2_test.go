// Álvaro Castellano Vela 2019/12/01
package main

import (
	"testing"
)

func TestMass12Fuel2(t *testing.T) {

	fuel := calculateFuel(12)
	if fuel != 2 {
		t.Errorf("calculateFuel(12) should return a 2, but returns %d", fuel)
	}
}

func TestMass1969Fuel966(t *testing.T) {

	fuel := calculateFuel(1969)
	if fuel != 966 {
		t.Errorf("calculateFuel(1969) should return a 966, but returns %d", fuel)
	}
}

func TestMass100756Fuel50346(t *testing.T) {

	fuel := calculateFuel(100756)
	if fuel != 50346 {
		t.Errorf("calculateFuel(100756) should return a 50346, but returns %d", fuel)
	}
}
