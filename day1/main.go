package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func getInput() ([]int, error) {
	var input []int

	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			return nil, err
		}

		row, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}

		input = append(input, row)
	}

	return input, nil
}

func part1(input []int) int {
	for _, i := range input {
		for _, j := range input {
			if i+j == 2020 {
				return i * j
			}
		}
	}

	return 0
}

func part2(input []int) int {
	for _, i := range input {
		for _, j := range input {
			for _, k := range input {
				if i+j+k == 2020 {
					return i * j * k
				}
			}
		}
	}

	return 0
}

func main() {
	input, err := getInput()
	if err != nil {
		log.Fatalf("Error getting input: %s", err)
	}

	var answer int

	//part1
	answer = part1(input)
	log.Printf("Part 1 Answer: %d", answer)

	//part2
	answer = part2(input)
	log.Printf("Part 2 Answer: %d", answer)
}
