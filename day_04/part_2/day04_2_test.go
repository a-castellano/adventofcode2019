// Álvaro Castellano Vela 2019/12/04
package main

import (
	"testing"
)

func TestNumberLargerThan6Digits(t *testing.T) {

	var result bool = meetCriteria(1234567)
	if result == true {
		t.Errorf("1234567 shouldn't meet criteria")

	}
}

func TestNumberShorterThan6Digits(t *testing.T) {

	var result bool = meetCriteria(12345)
	if result == true {
		t.Errorf("12345 shouldn't meet criteria")

	}
}

func TestNumberWithNoAdjacent01(t *testing.T) {

	var result bool = meetCriteria(654321)
	if result == true {
		t.Errorf("654321 shouldn't meet criteria")

	}
}

func TestNumberWithNoAdjacent02(t *testing.T) {

	var result bool = meetCriteria(654326)
	if result == true {
		t.Errorf("654326 shouldn't meet criteria")

	}
}

func TestNumberWithNoAdjacent03(t *testing.T) {

	var result bool = meetCriteria(634236)
	if result == true {
		t.Errorf("634236 shouldn't meet criteria")

	}
}

func TestNumberDecrease(t *testing.T) {

	var result bool = meetCriteria(654321)
	if result == true {
		t.Errorf("654321 shouldn't meet criteria")

	}
}

func TestNumberMeets01(t *testing.T) {

	var result bool = meetCriteria(135779)
	if result == false {
		t.Errorf("135779 should meet criteria")

	}
}

func TestNumberMeets02(t *testing.T) {

	var result bool = meetCriteria(123345)
	if result == false {
		t.Errorf("123345 should meet criteria")

	}
}

func TestNumberMeets03(t *testing.T) {

	var result bool = meetCriteria(112233)
	if result == false {
		t.Errorf("112233 should meet criteria")

	}
}

func TestNumberMeets04(t *testing.T) {

	var result bool = meetCriteria(111122)
	if result == false {
		t.Errorf("111122 should meet criteria")
	}
}

func TestNumberMeets05(t *testing.T) {

	var result bool = meetCriteria(111122)
	if result == false {
		t.Errorf("233344 should meet criteria")
	}
}

func TestNumberDoesNotMeet01(t *testing.T) {

	var result bool = meetCriteria(223450)
	if result == true {
		t.Errorf("223450 shouldn't meet criteria")

	}
}

func TestNumberDoesNotMeet02(t *testing.T) {

	var result bool = meetCriteria(123789)
	if result == true {
		t.Errorf("123789 shouldn't meet criteria")

	}
}

func TestNumberDoesNotMeet03(t *testing.T) {

	var result bool = meetCriteria(123444)
	if result == true {
		t.Errorf("123444 shouldn't meet criteria")

	}
}

func TestNumberDoesNotMeet05(t *testing.T) {

	var result bool = meetCriteria(124445)
	if result == true {
		t.Errorf("124445 shouldn't meet criteria")

	}
}
