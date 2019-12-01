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

func TestMass14Fuel2(t *testing.T) {

	fuel := calculateFuel(14)
	if fuel != 2 {
		t.Errorf("calculateFuel(14) should return a 2, but returns %d", fuel)
	}
}

func TestMass1969Fuel654(t *testing.T) {

	fuel := calculateFuel(1969)
	if fuel != 654 {
		t.Errorf("calculateFuel(1969) should return a 654, but returns %d", fuel)
	}
}

func TestMass100756Fuel33583(t *testing.T) {

	fuel := calculateFuel(100756)
	if fuel != 33583 {
		t.Errorf("calculateFuel(100756) should return a 33583, but returns %d", fuel)
	}
}
