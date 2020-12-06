package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var (
	sumPart1 int
	sumPart2 int
)

func checkFormAnswers() error {
	f, err := os.Open("input.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var buffer []string
	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			return err
		}

		line := scanner.Text()

		if len(line) == 0 {
			sumPart1 = sumPart1 + part1(buffer)
			sumPart2 = sumPart2 + part2(buffer)
			buffer = []string{}
			continue
		}

		buffer = append(buffer, scanner.Text())
	}

	if len(buffer) > 0 {
		sumPart1 = sumPart1 + part1(buffer)
		sumPart2 = sumPart2 + part2(buffer)
	}

	return nil
}

func part1(group []string) int {
	answers := map[string]struct{}{}
	for _, line := range group {
		for _, c := range line {
			answers[string(c)] = struct{}{}
		}
	}

	return len(answers)
}

func part2(group []string) int {
	if len(group) == 1 {
		return len(group[0])
	}

	answers := map[string]struct{}{}
	for i, line := range group {
		inBoth := map[string]struct{}{}
		for _, c := range line {
			_, ok := answers[string(c)]
			if ok || i == 0 {
				inBoth[string(c)] = struct{}{}
			}
		}
		answers = inBoth
	}

	return len(answers)
}

func main() {
	err := checkFormAnswers()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error parsing customs forms: %s", err))
	}

	log.Println(fmt.Sprintf("There are %d answered questions part 1", sumPart1))
	log.Println(fmt.Sprintf("There are %d answered questions part 2", sumPart2))
}
