package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
)

func getNumbers() (numbers []int, err error) {
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

		number, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, number)
	}

	return numbers, nil
}

func part1(numbers []int) (int, error) {
	diffs := map[int]int{
		1: 0,
		2: 0,
		3: 1,
	}

	last := 0
	for i := 0; i < len(numbers); i++ {
		diff := numbers[i] - last
		if diff > 3 {
			return 0, errors.New("difference was greater than 3")
		}
		diffs[diff]++
		last = numbers[i]
	}

	return diffs[1] * diffs[3], nil
}

func part2(numbers []int) int {
	numbers = append([]int{0}, numbers...)
	numbers = append(numbers, numbers[len(numbers)-1]+3)

	combos := 1
	setLens := map[int]int{}
	buffer := []int{}
	low := -1

	for i := 1; i < len(numbers); i++ {
		if numbers[i]-numbers[i-1] < 3 && numbers[i+1]-numbers[i] < 3 {
			if low == -1 {
				low = numbers[i-1]
			}
			buffer = append(buffer, numbers[i])
		} else {
			if len(buffer) > 0 {
				possible := int(math.Pow(2, float64(len(buffer))))
				if numbers[i]-low > 3 {
					possible = possible - 1
				}

				setLens[possible]++
				buffer = []int{}
				low = -1
			}
		}
	}

	for length, count := range setLens {
		combos *= int(math.Pow(float64(length), float64(count)))
	}

	return combos
}

func main() {
	numbers, err := getNumbers()
	sort.Ints(numbers)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error fetching adapters: %s", err))
	}

	product, err := part1(numbers)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error finding differences part 1: %s", err))
	}

	log.Println(fmt.Sprintf("The product of the difference is %d in part 1", product))

	sets := part2(numbers)
	log.Println(fmt.Sprintf("There are %d possible sets in part 2", sets))
}
