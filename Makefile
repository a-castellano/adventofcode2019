all: test

test: day_01 day_02

day_01: day_01_1 day_01_2
day_02: day_02_1

day_01_1: 
	cd day_01/part_1 && go test

day_01_2: 
	cd day_01/part_2 && go test

day_02_1: 
	cd day_02/part_1 && go test
