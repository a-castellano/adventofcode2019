// Álvaro Castellano Vela 2019/12/05
package main

import (
	"testing"
)

func TestDecodeInstruction01(t *testing.T) {

	var instruction int = 1002
	var opcode, mode1, mode2, mode3 int = decodeInstruction(instruction)

	if opcode != 2 {
		t.Errorf("Opcode whould be 2 instead of %d", opcode)
	}

	if mode1 != 0 {
		t.Errorf("mode1 whould be 0 instead of %d", mode1)
	}

	if mode2 != 1 {
		t.Errorf("mode2 whould be 1 instead of %d", mode2)
	}

	if mode3 != 0 {
		t.Errorf("mode3 whould be 0 instead of %d", mode3)
	}

}

func TestDecodeInstruction02(t *testing.T) {

	var instruction int = 10103
	var opcode, mode1, mode2, mode3 int = decodeInstruction(instruction)

	if opcode != 3 {
		t.Errorf("Opcode whould be 3 instead of %d", opcode)
	}

	if mode1 != 1 {
		t.Errorf("mode1 whould be 1 instead of %d", mode1)
	}

	if mode2 != 0 {
		t.Errorf("mode2 whould be 0 instead of %d", mode2)
	}

	if mode3 != 1 {
		t.Errorf("mode3 whould be 1 instead of %d", mode3)
	}

}
