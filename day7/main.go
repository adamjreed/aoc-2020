package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var bags map[string]map[string]int

func parseBags() error {
	bags = map[string]map[string]int{}
	f, err := os.Open("input.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			return err
		}

		line := scanner.Text()
		parts := strings.Split(line, " bags contain ")
		contains := strings.Split(parts[1][0:len(parts[1])-1], ", ")

		bags[parts[0]] = map[string]int{}
		for _, contained := range contains {
			if contained == "no other bags" {
				continue
			}

			howMany, err := strconv.Atoi(contained[0:1])
			if err != nil {
				return err
			}

			if howMany > 1 {
				bags[parts[0]][contained[2:len(contained)-5]] = howMany
			} else {
				bags[parts[0]][contained[2:len(contained)-4]] = howMany
			}

		}
	}

	return nil
}

func canContain(targetType string, bag map[string]int) bool {
	for bagType, _ := range bag {
		if bagType == targetType {
			return true
		}

		if canContain(targetType, bags[bagType]) {
			return true
		}
	}

	return false
}

func countContainedBags(bag map[string]int) int {
	var totalBags int
	for bagType, count := range bag {
		totalBags = totalBags + count + countContainedBags(bags[bagType])*count
	}

	return totalBags
}

func part1() int {
	var bagCombos int

	for bagType, contains := range bags {
		if bagType == "shiny gold" {
			continue
		}

		if canContain("shiny gold", contains) {
			bagCombos++
		}
	}

	return bagCombos
}

func part2() int {
	return countContainedBags(bags["shiny gold"])
}

func main() {
	err := parseBags()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error parsing bags: %s", err))
	}

	log.Println(fmt.Sprintf("There are %d possible bag combos in part 1", part1()))
	log.Println(fmt.Sprintf("There are %d total contained bags in part 2", part2()))
}
