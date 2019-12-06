all: test

test: day_01 day_02 day_04 day_05

day_01: day_01_1 day_01_2
day_02: day_02_1
day_04: day_04_1 day_04_2
day_05: day_05_1 day_05_2

day_01_1: 
	cd day_01/part_1 && go test

day_01_2: 
	cd day_01/part_2 && go test

day_02_1: 
	cd day_02/part_1 && go test

day_04_1: 
	cd day_04/part_1 && go test

day_04_2: 
	cd day_04/part_2 && go test

day_05_1: 
	cd day_05/part_1 && go test

day_05_2: 
	cd day_05/part_2 && go test
