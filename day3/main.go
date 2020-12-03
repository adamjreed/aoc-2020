package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func getInput() ([][]string, error) {
	var input [][]string

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

		input = append(input, strings.Split(scanner.Text(), ""))
	}

	return input, nil
}

func countTrees(right int, down int, slope [][]string) int {
	var (
		trees    int
		position int
	)

	slope = slope[down:]
	rowLength := len(slope[0])

	for i, row := range slope {
		if i%down != 0 {
			continue
		}

		position = (position + right) % rowLength
		if row[position] == "#" {
			trees++
		}
	}

	return trees
}

func part1(rows [][]string) int {
	return countTrees(3, 1, rows)
}

func part2(rows [][]string) int {
	product := 1

	paths := [][]int{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}

	for _, path := range paths {
		product = product * countTrees(path[0], path[1], rows)
	}

	return product
}

func main() {
	input, err := getInput()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Could not get input"))
	}

	log.Println(fmt.Sprintf("You encountered %d trees in part 1", part1(input)))
	log.Println(fmt.Sprintf("The product of trees you encountered in part 2 is %d", part2(input)))
}
